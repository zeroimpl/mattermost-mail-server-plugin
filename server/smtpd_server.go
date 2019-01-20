package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/chrj/smtpd"
	"github.com/mattermost/mattermost-server/model"
	"mime"
	//	"mime/multipart
	"errors"
	"net/textproto"
)

func (p *Plugin) runSMTPD() error {

	if !p.smtpdServerSetup {
		p.API.LogWarn("Launching smtpd server")
		p.smtpdServerSetup = true

		// run smtpd server in a separate goroutine
		go func(p *Plugin) {
			p.API.LogWarn("Runnning smtpd server")

			server := &smtpd.Server{

				//			HeloChecker: func(peer smtpd.Peer, name string) error {
				//				if !strings.HasPrefix(peer.Addr.String(), "42.42.42.42:") {
				//					return errors.New("Denied")
				//				}
				//				return nil
				//			},

				Handler: func(peer smtpd.Peer, env smtpd.Envelope) error {
					p.API.LogWarn("Received message", "sender", env.Sender)

					// Parse the email body as a MIME message, first extracting the headers and Content-Type
					reader := bytes.NewReader(env.Data)
					tpReader := textproto.NewReader(bufio.NewReader(reader))
					mimehdr, err := tpReader.ReadMIMEHeader()
					if err != nil {
						p.API.LogWarn("Mime Header error", "err", err)
						return nil
					}
					// p.API.LogWarn("Mime Header", "mime", fmt.Sprintf("%s", mimehdr))

					ct, ok := mimehdr["Content-Type"]
					if !ok {
						p.API.LogWarn("No content type")
						return errors.New("No Content-Type header found")
					}

					mediatype, mtparams, err := mime.ParseMediaType(ct[0])
					if err != nil {
						p.API.LogWarn("Media type parse error", "err", err)
						return err
					}

					p.API.LogWarn("Media type", "media type", mediatype)
					p.API.LogWarn("Media type params", "media type params", fmt.Sprintf("%s", mtparams))

					txt := ""
					if mediatype == "text/plain" {
						for {
							s, err := tpReader.ReadLine()
							if err != nil {
								break
							}
							txt += s + "\n"
						}
						// Try and parse multipart message. TODO: This code doesn't work yet
						//	} else if mediatype == "multipart/alternative" {
						//		boundary := mtparams["boundary"]
						//		mreader := multipart.NewReader(reader, boundary)
						//		// look for a text/plain section
						//		for {
						//			part, err := mreader.NextPart()
						//			if err != nil || part == nil {
						//				break
						//			}
						//			buf := new(bytes.Buffer)
						//			buf.ReadFrom(part)
						//			txt += buf.String() + "\n\n\n"
						//			break
						//		}
					} else {
						txt = fmt.Sprintf("Email uses an unsupported media type, %s, raw data is as follows:\n%s", mediatype, env.Data)
					}

					if _, err := p.API.CreatePost(&model.Post{
						// hardcoded values. These should be computed from the receipient address
						UserId:    "5yx5ebr97ffh9kbcsfgf7o87ry",
						ChannelId: "wizz31n4kbrfmnjpnw9ra7ygtw",
						Message:   fmt.Sprintf("Email received: From=%s\nTo=%s\n%s", env.Sender, env.Recipients, txt),
					}); err != nil {
						p.API.LogError(
							"failed to post message",
						)
						return err
					}

					return nil
				},
			}

			// Listen on localhost on port 10025. TODO: This should be configurable and default to port 0.0.0.0:25
			server.ListenAndServe("127.0.0.1:10025")
		}(p)

	} else {
		p.API.LogWarn("smtpd server was already running")
	}

	return nil
}
