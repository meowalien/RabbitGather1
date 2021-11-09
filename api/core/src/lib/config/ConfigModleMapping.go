package config

import "github.com/spf13/viper"

// JsonConfigModleMapping mapping the given json config file to the given struct.
// Using `mapstructure:"<name>"` tag to customize the field name
func JsonConfigModleMapping(stc interface{}, filePath string) error {
	viper.SetConfigFile(filePath)
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(stc)
	if err != nil {
		return err
	}
	return nil
}
