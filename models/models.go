package models

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
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