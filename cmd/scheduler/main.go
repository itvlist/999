package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/ajvb/kala/client"
	"github.com/ajvb/kala/job"
	"github.com/schollz/closestmatch"
	"github.com/urfave/cli"
)

const dateLayout = "20060102150405 -0700"

type TV struct {
	XMLName  xml.Name  `xml:"tv"`
	Channels []Channel `xml:"channel"`
	Programs []Program `xml:"programme"`
}

type Channel struct {
	XMLName  xml.Name `xml:"channel"`
	Id       string   `xml:"id,attr"`
	Programs []Program
}

func (c *Channel) ParseProgram(ps []Program, t time.Time) {
	for _, p := range ps {
		if p.Channel == c.Id && p.StartTime().YearDay() == t.YearDay() {
			c.Programs = append(c.Programs, p)
		}
	}
}

type Program struct {
	XMLName xml.Name `xml:"programme"`
	Start   string   `xml:"start,attr"`
	Stop    string   `xml:"stop,attr"`
	Channel string   `xml:"channel,attr"`
	Title   string   `xml:"title"`
}

func (p Program) StartTime() time.Time {
	t, _ := time.Parse(dateLayout, p.Start)
	return t
}

func (p Program) StopTime() time.Time {
	t, _ := time.Parse(dateLayout, p.Stop)
	return t
}
func (p Program) StartDay() string {
	t := p.StartTime()
	return fmt.Sprintf("%02d/%02d/%02d", t.Day(), t.Month(), t.Year())

}
func (p Program) StartHour() string {
	t := p.StartTime()
	return fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
}

func (p Program) String() string {
	return fmt.Sprintf("%s %s: %s", p.StartDay(), p.StartHour(), p.Title)
}

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "channel",
		},
		cli.StringFlag{
			Name: "date",
		},
		cli.StringFlag{
			Name: "title",
		},
	}

	app.Name = "boom"
	app.Usage = "make an explosive entrance"
	app.Action = func(ctx *cli.Context) error {
		p, err := parseEPG(ctx)

		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		err = scheduleRecord(p)

		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		return nil
	}

	app.Run(os.Args)
}

func scheduleRecord(p Program) error {
	c := client.New("http://127.0.0.1:8000")

	body := &job.Job{
		Schedule: fmt.Sprintf("R0/%s/PT1S", time.Now().Add(10*time.Second).Format(time.RFC3339)),
		Name:     fmt.Sprintf("record_%04d", rand.Intn(5000)),
		Command:  fmt.Sprintf("iptv-record --channel \"%s\" --starttime \"%s\" --stoptime \"%s\"", p.Channel, p.StartTime().Format(time.RFC3339), p.StopTime().Format(time.RFC3339)),
	}
	id, err := c.CreateJob(body)
	log.Println("Job Created: ", id)
	if err != nil {
		return err
	}

	return nil
}
func parseEPG(ctx *cli.Context) (Program, error) {

	xmlFile, err := os.Open("xmltv.xml")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened users.xml")
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var tv TV
	xml.Unmarshal(byteValue, &tv)
	title := ctx.String("title")
	fmt.Printf("Looking for %s on %s the %s\n", title, ctx.String("channel"), ctx.String("date"))
	t, _ := time.Parse("02-01-2006", ctx.String("date"))
	for _, c := range tv.Channels {
		if c.Id == ctx.String("channel") {
			c.ParseProgram(tv.Programs, t)
			ps := c.Programs
			if title != "" {
				ps = c.ClosestMatch(title)
			}
			p, err := askForChoice(ps)
			if err != nil {
				return Program{}, err
			}

			return p, nil
		}
	}

	return Program{}, nil
}

func (c *Channel) ClosestMatch(title string) []Program {
	ps := []string{}
	for _, p := range c.Programs {
		ps = append(ps, p.Title)
	}
	bagSizes := []int{3}
	cm := closestmatch.New(ps, bagSizes)
	r := cm.ClosestN(title, 2)

	result := []Program{}

	for _, x := range r {

		for _, p := range c.Programs {
			if p.Title == x {
				result = append(result, p)
			}
		}
	}

	return result
}

func askForChoice(programs []Program) (Program, error) {
	for true {
		for i, v := range programs {
			fmt.Printf("%d) %s\n", i, v)
		}
		fmt.Printf("Choice: ")
		var choice int
		_, err := fmt.Scan(&choice)

		if err != nil || choice >= len(programs) {
			fmt.Println("Choice not correct")
		} else {
			return programs[choice], nil
		}
	}

	return Program{}, nil
}