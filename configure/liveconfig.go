package configure

/*
{
	[
	{
	"application":"live",
	"live":"on",
	"hls":"on",
	"static_push":["rtmp://xx/live"]
	}
	]
}
*/
type Application struct {
	Appname     string
	Liveon      string
	Hlson       string
	Static_push []string
}

type ServerCfg struct {
	Server []Application
}

var RtmpServercfg ServerCfg

func LoadConfig() error {
	RtmpServercfg.Server = append(RtmpServercfg.Server, Application{
		Appname: "live",
		Liveon:  "on",
		Hlson:   "on",
	})

	return nil
}

func CheckAppName(appname string) bool {
	for _, app := range RtmpServercfg.Server {
		if (app.Appname == appname) && (app.Liveon == "on") {
			return true
		}
	}
	return false
}

func GetStaticPushUrlList(appname string) ([]string, bool) {
	for _, app := range RtmpServercfg.Server {
		if (app.Appname == appname) && (app.Liveon == "on") {
			if len(app.Static_push) > 0 {
				return app.Static_push, true
			} else {
				return nil, false
			}
		}

	}
	return nil, false
}
