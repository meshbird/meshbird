package cliconfig

import (
	"reflect"
	"strconv"
	"strings"
	"os"

	"github.com/serenize/snaker"
	"github.com/urfave/cli"
)

//Fill func get config struct and envPrefix to collect Urfave Cli slice of cli.Flag
func Fill(config interface{}, envPrefix string) []cli.Flag {
	configValue := reflect.Indirect(reflect.ValueOf(config))
	var flags []cli.Flag
	for i := 0; i < configValue.NumField(); i++ {
		fieldValue := configValue.Field(i)
		fieldType := configValue.Type().Field(i)
		name := snaker.CamelToSnake(fieldType.Name)
		flagName := fieldType.Tag.Get("flag")
		if flagName == "" {
			flagName = name
		}
		envName := fieldType.Tag.Get("env")
		if envName == "" {
			envName = strings.ToUpper(flagName)
		}
		envName = envPrefix + envName
		switch fieldType.Type.Kind() {
		case reflect.String:
			flag := cli.StringFlag{
				Name:        flagName,
				EnvVar:      envName,
				Destination: fieldValue.Addr().Interface().(*string),
				Value:       fieldType.Tag.Get("default"),
			}
			flags = append(flags, flag)
		case reflect.Int:
			flag := cli.IntFlag{
				Name:        flagName,
				EnvVar:      envName,
				Destination: fieldValue.Addr().Interface().(*int),
				Value:       intFromString(fieldType.Tag.Get("default")),
			}
			flags = append(flags, flag)
		case reflect.Slice:
			if fieldType.Type.Elem().Kind() == reflect.String {
				values := strings.Split(fieldType.Tag.Get("default"), ",")
				values2 := cli.StringSlice(values)
				fieldValue.Set(reflect.ValueOf(values))
				flag := cli.StringSliceFlag {
					Name:        flagName,
					EnvVar:      envName,
//					Destination: fieldValue.Addr().Interface().(*[]string),
					Value:       &values2,
				}
				flags = append(flags, flag)
			}
		}
	}
	return flags
}

func intFromString(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}

//FillAndRun function for simpler fill config and run cli app
func FillAndRun(config interface{}, appName, envPrefix string, handler func(*cli.Context) error) error {
	app := cli.NewApp()
	app.Name = appName
	app.Flags = Fill(&config, envPrefix)
	app.Action = handler
	return app.Run(os.Args)
}