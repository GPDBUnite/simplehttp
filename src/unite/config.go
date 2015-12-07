package unite

import (
       "fmt"    
       "io"     
       "os"
       "strings"
       "bufio"
)

type Config map[string]string
type ConfigSet map[string]Config

func ParseConfig(filename string) ConfigSet {
	ret := make(ConfigSet)
	c := make(Config)
	configname := "default"

	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer f.Close()
	r := bufio.NewReader(f)
    var line string 
	for err == nil {
		line, err = r.ReadString('\n')
        // fmt.Printf("%s",line)
        if strings.HasPrefix(line, "#") {
           // comments, skip
           continue
        }
		if strings.HasPrefix(line, "[") {
			ret[configname] = c
			i := strings.LastIndex(line, "]")
			configname = line[1:i]
			c = make(Config)
		} else {
			s := strings.SplitN(line, "=", 2)
            if len(s) < 2 {
               // bad format, skip? quit?
               // fmt.Printf("skip %s\n", line)
               continue
            }
			c[strings.TrimRight(s[0], "\t ")] = strings.TrimRight(s[1], "\r\n ")
		}
	}
	if err != io.EOF {
		fmt.Println(err)
		return ret
	} else {
		if strings.Count(line, "") > 1 {
			if strings.HasPrefix(line, "[") {
				ret[configname] = c
				i := strings.LastIndex(line, "]")
				configname = line[1:i]
				c = make(Config)
			} else {
				s := strings.SplitN(line, "=", 2)
				c[s[0]] = strings.TrimRight(s[1], "\r\n ")
			}
		}
		ret[configname] = c
	}

	return ret
}

func (csets ConfigSet) GetConfig(sec string, item string, defvalue string) string {
	if s, ok := csets[sec]; ok {
		if v, ok := s[item]; ok {
			return v
		}
	}
	return defvalue
}
