package clickcounteriframe

import (
	"log"
	"net/http"
	"text/template"
)

type IframeController struct {
	domainRepo DomainRepositoryInterface
}

var iframeTpl = `<!DOCTYPE html>
<!--


            _|_|          _|_|
            _|_|          _|_|
            _|_|
            _|_|
            _|_|_|_|      _|_|  _|_|    _|_|
            _|_|_|_|_|    _|_|  _|_|    _|_|
            _|_|    _|_|  _|_|  _|_|    _|_|
            _|_|    _|_|  _|_|  _|_|    _|_|
            _|_|    _|_|  _|_|  _|_|    _|_|
      _|_|  _|_|    _|_|  _|_|    _|_|_|_|
      _|_|  _|_|    _|_|  _|_|      _|_|

      .hiv domains – The digital Red Ribbon

                  click4life.hiv

-->
<html>
<head>
    <title>{{.Name}}</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style type="text/css">
        html, body {
            margin: 0;
            padding: 0;
            width: 100%;                                                                                                                                                                                                                     
            height: 100%;                                                                                                                                                                                                                    
            overflow: hidden;                                                                                                                                                                                                                
        }                                                                                                                                                                                                                                    
                                                                                                                                                                                                                                             
        #clickcounter-target-iframe {                                                                                                                                                                                                        
            border: 0;                                                                                                                                                                                                                       
            width: 100%;                                                                                                                                                                                                                     
            height: 100%;                                                                                                                                                                                                                    
            margin: 0;                                                                                                                                                                                                                       
            padding: 0;                                                                                                                                                                                                                      
        }                                                                                                                                                                                                                                    
    </style>                                                                                                                                                                                                                                 
</head>                                                                                                                                                                                                                                      
<body>                                                                                                                                                                                                                                       
<iframe src="{{.Redirect}}" width="100%" height="100%" id="clickcounter-target-iframe"></iframe>                                                                                                                                           
<script src="//dothiv-registry.appspot.com/static/clickcounter.min.js" type="text/javascript"></script>                                                                                                                                      
</body>                                                                                                                                                                                                                                      
</html>                                                                      
`

func (c *IframeController) IframeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(400)
		return
	}
	w.Header().Add("Content-Type", "text/html; charset=UTF-8")
	w.Header().Add("Cache-Control", "public, s-maxage=1800")

	domain := new(Domain)
	domain.Name = "thjnk.hiv"
	domain.Redirect = "http://thjnk.de"

	t := template.Must(template.New("iframe").Parse(iframeTpl))
	err := t.Execute(w, domain)
	if err != nil {
		log.Fatalln("failed to parse template:", err)
	}

	// w.Write([]byte(strings.Split(r.Host, ":")[0]))
}