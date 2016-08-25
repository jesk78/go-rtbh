package proto

import (
	"github.com/r3boot/go-rtbh/config"
	"github.com/r3boot/rlib/sys"
	"os"
	"text/template"
)

type Bird struct {
}

type BirdPeer struct {
	Name    string
	Address string
	PeerAs  string
}

type BirdTemplate struct {
	Asnum     string
	RouterId  string
	Whitelist []string
	Blacklist []string
	Peers     []*BirdPeer
}

type BirdStatusProtocolEntry struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Table    string `json:"table"`
	State    string `json:"state"`
	Since    string `json:"since"`
	Info     string `json:"info"`
}

type BirdStatusOutput struct {
	Protocols []BirdStatusProtocolEntry `json:"protocols"`
}

func LoadTemplate() []byte {
	var fs os.FileInfo
	var fd *os.File
	var err error
	var size int
	var data []byte

	if fs, err = os.Stat(config.TMPL_BIRD); err != nil {
		Log.Fatal("[Bird.LoadTemplate]: Failed to load " + config.TMPL_BIRD)
		return []byte{}
	}

	data = make([]byte, fs.Size())

	if fd, err = os.Open(config.TMPL_BIRD); err != nil {
		Log.Fatal("[Bird.LoadTemplate]: Failed to open " + config.TMPL_BIRD)
		return []byte{}
	}

	if size, err = fd.Read(data); err != nil {
		Log.Fatal("[Bird.LoadTemplate]: Failed to read " + config.TMPL_BIRD)
		return []byte{}
	}

	if int64(size) != fs.Size() {
		Log.Fatal("[Bird.LoadTemplate]: Size mismatch on " + config.TMPL_BIRD)
		return []byte{}
	}

	return data
}

func (b *Bird) Command(cmdline string) (output []string, err error) {
	output, _, err = sys.Run("birdc", cmdline)
	return
}

func (b *Bird) WriteTo(cfgfile string) bool {
	var err error

	if _, err = os.Stat(cfgfile); err != nil {
		Log.Warning("[Bird.WriteTo]: " + cfgfile + " does not exist")
		return false
	}

	return true
}

func (b *Bird) ExportPrefixes(wl []string, bl []string) bool {
	var fd *os.File
	var err error
	var template_data []byte
	// var cfg_data bytes.Buffer

	template_data = LoadTemplate()
	if len(template_data) == 0 {
		Log.Warning("[Bird.ExportPrefixes]: Found empty template")
		return false
	}

	tmpl, err := template.New("config").Parse(string(template_data))
	if err != nil {
		Log.Warning("[Bird.ExportPrefixes]: Failed to prepare template")
		Log.Warning(err)
		return false
	}

	data := &BirdTemplate{
		Asnum:     Config.BGP.Asnum,
		RouterId:  Config.BGP.RouterId,
		Whitelist: wl,
		Blacklist: bl,
	}

	for _, peer := range Config.BGP.Peers {
		router := &BirdPeer{
			Name:    peer.Name,
			Address: peer.Address,
			PeerAs:  peer.Asnum,
		}

		data.Peers = append(data.Peers, router)
	}

	if _, err = os.Stat(Config.BGP.ConfigFile); err == nil {
		if err = os.Remove(Config.BGP.ConfigFile); err != nil {
			Log.Warning("[Bird.ExportPrefixes]: Failed to unlink " + Config.BGP.ConfigFile)
			return false
		}
	}

	fd, err = os.OpenFile(Config.BGP.ConfigFile, (os.O_CREATE | os.O_WRONLY), 0644)
	if err != nil {
		Log.Warning("[Bird.ExportPrefixes]: Failed to open " + Config.BGP.ConfigFile)
		return false
	}
	defer fd.Close()

	if tmpl.Execute(fd, data); err != nil {
		Log.Warning("[Bird.ExportPrefixes]: Failed to execute template: " + err.Error())
		return false
	}

	Log.Debug("[Bird.ExportPrefixes]: Wrote " + Config.BGP.ConfigFile)

	if _, err = b.Command("configure"); err != nil {
		Log.Warning("[Bird.ExportPrefixes]: Failed to reload bird configuration")
		return false
	}
	Log.Debug("[Bird.ExportPrefixes]: Reloaded bird configuration")

	return true
}

func NewBirdClient() (b *Bird) {
	b = &Bird{}

	return
}
