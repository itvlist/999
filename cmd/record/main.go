package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/urfave/cli"
)

type Channel struct {
	TVGID      string `json:"tvgid"`
	TVGName    string `json:"tvg_name"`
	TVGLogo    string `json:"tvg_logo"`
	URL        string `json:"url"`
	GroupTitle string `json:"group_title"`
	Name       string `json:"name"`
}

type Group struct {
	Name string `json:"name"`
}

func (c Channel) String() string {
	return fmt.Sprintf("tvg-id: %s, tvg-name: %s, tvg-logo: %s, group-title: %s, name: %s, url: %s", c.TVGID, c.TVGName, c.TVGLogo, c.GroupTitle, c.Name, c.URL)
}

func Parse(fileName string) ([]Channel, error) {
	var regex = regexp.MustCompile(`^#EXTINF:-1 tvg-id="(?P<TVGID>[^"]*)" tvg-name="(?P<TVGName>[^"]*)" tvg-logo="(?P<TVGLogo>[^"]*)" group-title="(?P<GroupTitle>[^"]*)"`)
	f, err := os.Open(fileName)
	if err != nil {
		return nil, errors.New("Unable to open playlist file")
	}
	defer f.Close()

	onFirstLine := true
	scanner := bufio.NewScanner(f)

	channels := []Channel{}
	for scanner.Scan() {
		line := scanner.Text()
		if onFirstLine && !strings.HasPrefix(line, "#EXTM3U") {
			return nil, errors.New("Invalid m3u file format. Expected #EXTM3U file header")
		}

		onFirstLine = false

		if strings.HasPrefix(line, "#EXTINF") {
			s := reSubMatchMap(regex, line)
			channel := Channel{}
			channel.TVGID = s["TVGID"]
			channel.TVGName = s["TVGName"]
			channel.TVGLogo = s["TVGLogo"]
			channel.GroupTitle = s["GroupTitle"]
			channel.Name = s["TVGName"]
			scanner.Scan()
			channel.URL = scanner.Text()
			channels = append(channels, channel)
		}
	}

	return channels, nil
}

func GetGroups(channels []Channel) []Group {
	groups := make(map[string]interface{})

	for _, c := range channels {
		groups[c.GroupTitle] = true
	}

	keys := make([]Group, 0, len(groups))
	for k := range groups {
		keys = append(keys, Group{Name: k})
	}

	return keys
}

func GetChannels(channels []Channel, groupName string) []Channel {
	c := GetChannelByGroup(channels)

	return c[groupName]
}

func GetChannelByGroup(channels []Channel) map[string][]Channel {

	groups := make(map[string][]Channel)
	for _, c := range channels {
		if _, ok := groups[c.GroupTitle]; !ok {
			groups[c.GroupTitle] = []Channel{}
		}
		groups[c.GroupTitle] = append(groups[c.GroupTitle], c)
	}
	return groups
}

func GetChannel(channels []Channel, id string) (Channel, bool) {
	cs := []Channel{}
	for _, c := range channels {
		if c.TVGID == id {
			cs = append(cs, c)
		}
	}
	if len(cs) == 0 {

		return Channel{}, false
	}

	for _, c := range cs {
		if strings.Contains(c.Name, "HD") {
			return c, true
		}
	}
	return cs[0], true
}

func reSubMatchMap(r *regexp.Regexp, str string) map[string]string {
	match := r.FindStringSubmatch(str)
	subMatchMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			subMatchMap[name] = match[i]
		}
	}
	return subMatchMap
}
func main() {
	f, err := os.OpenFile("/home/jgsqware/iptv-record.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "channel",
		},
		cli.StringFlag{
			Name: "starttime",
		},
		cli.StringFlag{
			Name: "stoptime",
		},
	}

	app.Name = "iptv-record"
	app.Usage = "Record flux from ts file with duration"
	app.Action = func(ctx *cli.Context) error {
		channels, err := Parse("/home/jgsqware/go/src/github.com/jgsqware/iptv-record/playlist.m3u")
		if err != nil {
			log.Fatal(err)
		}
		channel := ctx.String("channel")
		c, ok := GetChannel(channels, channel)

		if !ok {
			return cli.NewExitError(fmt.Errorf("Channel not found: %s", channel), 1)
		}

		log.Printf("recording %s from %s to %s\n", c.Name, ctx.String("starttime"), ctx.String("stoptime"))

		starttime, err := time.Parse(time.RFC3339, ctx.String("starttime"))
		if err != nil {
			log.Fatal(err)
		}

		stoptime, err := time.Parse(time.RFC3339, ctx.String("stoptime"))
		if err != nil {
			log.Fatal(err)
		}

		delta := starttime.Sub(time.Now())

		log.Println("Will start in", fmtDuration(delta))
		time.Sleep(delta)
		delta = stoptime.Sub(starttime)

		log.Println("Will record for:", fmtDuration(delta))
		cmd := exec.Command("ffmpeg", "-i", c.URL, fmt.Sprintf("%s-%s.mkv", strings.Replace(channel, " ", "", -1), starttime.Format("02012006-1504")))
		log.Printf("Running ffmpeg...")
		err = cmd.Start()
		if err != nil {
			log.Fatal(err)
		}

		timer := time.NewTimer(delta)
		go func(timer *time.Timer, cmd *exec.Cmd) {
			for _ = range timer.C {
				log.Println("Kill ffmpeg")
				err := cmd.Process.Signal(os.Kill)
				if err != nil {
					log.Fatal(err)
				}
			}
		}(timer, cmd)

		log.Println("Recording...")
		cmd.Wait()
		log.Println("Command finished")
		return nil
	}

	app.Run(os.Args)

}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d hours %02d minutes", h, m)
}