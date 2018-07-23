/*
 * The MIT License (MIT)
 *
 * Copyright (c) 2018 Yu Jing <yu@argcv.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package config

import (
	"fmt"
	"github.com/argcv/webeh/log"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Option struct {
	Project              string
	Path                 string
	FileMustExists       bool
	DefaultPath          string
	DefaultConfigureName string
}

func (c *Option) With(rhs *Option) *Option {
	if rhs.Project != "" {
		c.Project = rhs.Project
	}
	if rhs.Path != "" {
		c.Path = rhs.Path
	}
	if rhs.FileMustExists {
		c.FileMustExists = true
	}
	if rhs.DefaultPath != "" {
		c.DefaultPath = rhs.DefaultPath
	}
	if rhs.DefaultConfigureName != "" {
		c.DefaultConfigureName = rhs.DefaultConfigureName
	}
	return c
}

func (c *Option) GetDefaultPath() string {
	if c.DefaultPath != "" {
		return c.DefaultPath
	} else {
		return c.Path
	}
}

// this function will search and load configurations
func LoadConfig(options ...Option) (err error) {
	option := Option{}

	for _, opt := range options {
		option.With(&opt)
	}

	project := option.Project

	if project == "" {
		return errors.New("Required parameter missing: project")
	}

	viper.SetConfigName(project)
	viper.SetEnvPrefix(project)

	if option.Path != "" {
		viper.SetConfigFile(option.Path)
	} else {
		viper.AddConfigPath(".")  // current folder
		viper.AddConfigPath("..") // parent folder
		viper.AddConfigPath("$HOME/")
		viper.AddConfigPath(fmt.Sprintf("$HOME/.%s/", project))
		viper.AddConfigPath("/etc/")
		viper.AddConfigPath(fmt.Sprintf("/etc/%s/", project))
		if conf := os.Getenv(fmt.Sprintf("%s_CFG", strings.ToUpper(project))); conf != "" {
			viper.SetConfigFile(conf)
		}
	}

	readAndTestConfig := func() (string, error) {
		err = viper.ReadInConfig()
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok && err != nil {
			return "", err
		}
		if conf := viper.ConfigFileUsed(); conf != "" {
			wd, _ := os.Getwd()
			if rel, _ := filepath.Rel(wd, conf); rel != "" && strings.Count(rel, "..") < 3 {
				conf = rel
			}
			log.Infof("using config file: %s", conf)
			return conf, nil
		} else {
			return "", nil
		}
	}

	if conf, e := readAndTestConfig(); conf != "" {
		return nil
	} else if e != nil {
		return e
	} else if option.FileMustExists {
		defaultConfigPath := option.GetDefaultPath()
		defaultConfigName := option.DefaultConfigureName
		if defaultConfigName == "" {
			defaultConfigName = fmt.Sprintf("%s.yml", project)
		}
		if defaultConfigPath == "" {
			msg := "No configure file: default path Not Assigned"
			log.Warnf(msg)
			return errors.New(msg)
		}
		if e := os.MkdirAll(defaultConfigPath, 0700); e != nil {
			return e
		}
		defaultConfigPathFile := path.Join(defaultConfigPath, defaultConfigName)
		log.Infof("configure file NOT found, using default file: %v", defaultConfigPathFile)
		if e := viper.WriteConfigAs(defaultConfigPathFile); e != nil {
			return e
		}
		viper.SetConfigFile(defaultConfigPathFile)

		if conf, e := readAndTestConfig(); conf != "" {
			return nil
		} else if e != nil {
			return e
		} else {
			msg := "No configure file"
			log.Warnf(msg)
			return errors.New(msg)
		}
	} else {
		return nil
	}
}
