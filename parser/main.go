package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"

	bazooka "github.com/bazooka-ci/bazooka/commons"
	bzklog "github.com/bazooka-ci/bazooka/commons/logs"
)

const (
	SourceFolder = "/bazooka"
	OutputFolder = "/bazooka-output"
	MetaFolder   = "/meta"
	RubyLang     = "ruby"
)

func init() {
	log.SetFormatter(&bzklog.BzkFormatter{})
	err := bazooka.LoadCryptoKeyFromFile("/bazooka-cryptokey")
	if err != nil {
		log.Fatal(err)
	}
}

type Configruby struct {
	Base         bazooka.Config `yaml:",inline"`
	RubyVersions []string       `yaml:"ruby,omitempty"`
}

func main() {
	file, err := bazooka.ResolveConfigFile(SourceFolder)
	if err != nil {
		log.Fatal(err)
	}

	conf := &Configruby{}
	err = bazooka.Parse(file, conf)
	if err != nil {
		log.Fatal(err)
	}

	if len(conf.Base.Script) == 0 {
		log.Fatal("Ruby builds should define a script value in the build descriptor")
	}

	versions := conf.RubyVersions
	images := conf.Base.Image

	if len(versions) == 0 && len(images) == 0 {
		versions = []string{"2.2"}
	}

	for i, version := range versions {
		if err := manageRubyVersion(fmt.Sprintf("0%d", i), conf, version, ""); err != nil {
			log.Fatal(err)
		}
	}
	for i, image := range images {
		if err := manageRubyVersion(fmt.Sprintf("1%d", i), conf, "", image); err != nil {
			log.Fatal(err)
		}
	}
}

func manageRubyVersion(counter string, conf *Configruby, version, image string) error {
	conf.RubyVersions = nil
	conf.Base.Image = nil

	meta := map[string]string{}
	if len(version) > 0 {
		var err error
		image, err = resolveRubyImage(version)
		if err != nil {
			return err
		}
		meta[RubyLang] = version
	} else {
		meta["image"] = image
	}
	conf.Base.FromImage = image

	if err := bazooka.Flush(meta, fmt.Sprintf("%s/%s", MetaFolder, counter)); err != nil {
		return err
	}
	return bazooka.Flush(conf, fmt.Sprintf("%s/.bazooka.%s.yml", OutputFolder, counter))
}

func resolveRubyImage(version string) (string, error) {
	//TODO extract this from db
	rubyMap := map[string]string{
		"2.0":         "bazooka/runner-ruby:2.0",
		"2.1":         "bazooka/runner-ruby:2.1",
		"2.2":         "bazooka/runner-ruby:2.2",
		"jruby1.7.20": "bazooka/runner-ruby:jruby1.7.20",
	}
	if val, ok := rubyMap[version]; ok {
		return val, nil
	}
	return "", fmt.Errorf("Unable to find Bazooka Docker Image for ruby Runnner %s", version)
}
