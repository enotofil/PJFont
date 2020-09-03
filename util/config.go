package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	configName = "config.json"
)

// UnicodeRange : First <= code <= Last
type UnicodeRange struct {
	Name           string
	FirstCodePoint uint
	LastCodePoint  uint
}

// Config data
type Config struct {
	TTFontsPath   string
	OutFontsPath  string
	SampleStrings []string
	UnicodeRanges []UnicodeRange
}

// LoadConfig or create new if error opening/parsing
func LoadConfig() *Config {
	var config Config
	f, err := os.Open(configName)
	if err != nil {
		fmt.Printf("cant open '%s': %s\n", configName, err)
		return createDefault()
	}
	defer f.Close()
	byteValue, _ := ioutil.ReadAll(f)
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		fmt.Printf("cant parse '%s': %s\n", configName, err)
		return createDefault()
	}
	fmt.Printf("'%s' loaded and parsed ok\n", configName)
	return &config
}

func createDefault() *Config {
	c := Config{
		TTFontsPath:  "./fonts",
		OutFontsPath: "./out",
		SampleStrings: []string{
			"Jackdaws love my big sphinx of quartz.",
			"В чащах юга жил бы цитрус? Да, но фальшивый!",
		},
		UnicodeRanges: []UnicodeRange{
			{
				Name:           "Basic",
				FirstCodePoint: 0x20,
				LastCodePoint:  0x7F,
			},
			{
				Name:           "Russian",
				FirstCodePoint: 0x410,
				LastCodePoint:  0x44F,
			},
		},
	}
	c.save()
	return &c
}

func (conf *Config) save() error {
	jsonBytes, _ := json.MarshalIndent(conf, "", "\t")
	return ioutil.WriteFile(configName, jsonBytes, os.ModePerm)
}

// CollectRunes from all ranges
func (conf *Config) CollectRunes() []rune {
	runes := []rune{}
	for _, ran := range conf.UnicodeRanges {
		for i := ran.FirstCodePoint; i <= ran.LastCodePoint; i++ {
			runes = append(runes, rune(i))
		}
	}
	return runes
}
