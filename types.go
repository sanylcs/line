package line

import (
	"errors"
	"fmt"
	"html"
	"strconv"
	"strings"
)

type Profile struct {
	DisplayName   string `json:"displayName"`
	Mid           string `json:"mid"`
	PictureUrl    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
}

type Profiles struct {
	Contacts *Profile `json:"contacts"`
	Count    int      `json:"count"`
	Total    int      `json:"total"`
	Start    int      `json:"start"`
	Display  int      `json:"display"`
}

type Group struct {
	Mid        string `json:"mid"`
	Name       string `json:"name"`
	PictureUrl string `json:"pictureUrl"`
}

type Groups struct {
	Groups  *Group `json:"groups"`
	Count   int    `json:"count"`
	Total   int    `json:"total"`
	Start   int    `json:"start"`
	Display int    `json:"display"`
}

type Placeholder map[string]string
type EscapePlaceholder map[string]string

func (e EscapePlaceholder) MarshalerJSON() ([]byte, error) {
	return placeholderMarshaler(e, true)
}

func (p Placeholder) MarshalerJSON() ([]byte, error) {
	return placeholderMarshaler(p, false)
}

func placeholderMarshaler(p map[string]string, htmlEscape bool) ([]byte, error) {
	if p == nil || len(p) == 0 {
		return nil, errors.New("can not marshal empty map")
	}
	a := make([]string, len(p))
	i := 0
	for k, v := range p {
		v = strconv.Quote(v)
		if htmlEscape {
			v = html.EscapeString(v)
		}
		a[i] = fmt.Sprint(strconv.Quote(k), ":", v)
		i++
	}
	s := strings.Join(a, `,`)
	return []byte(fmt.Sprint("{", s, "}")), nil
}

type LinkMsgContent struct {
	TemplateId     string            `json:"templateId"`
	PreviewUrl     string            `json:"previewUrl,omitempty"`
	TextParams     EscapePlaceholder `json:"textParams,omitempty"`
	SubTextParams  Placeholder       `json:"subTextParams,omitempty"`
	AltTextParams  Placeholder       `json:"altTextParams,omitempty"`
	LinkTextParams Placeholder       `json:"linkTextParams,omitempty"`
	ALinkUriParams Placeholder       `json:"aLinkTextParams,omitempty"`
	ILinkUriParams Placeholder       `json:"iLinkUriParams,omitempty"`
	LinkUriParams  Placeholder       `json:"linkUriParams,omitempty"`
}

type LinkMsg struct {
	To        []string        `json:"to"`
	ToChannel int             `json:"toChannel"`
	EventType string          `json:"eventType"`
	Content   *LinkMsgContent `json:"content"`
}

type PostResponse struct {
	Version   string   `json:"version"`
	Timestamp int64    `json:"timestamp"`
	MessageId string   `json:"messageId"`
	Failed    []string `json:"failed"`
}

type TimelineTemplate struct {
	DyamicObjs Placeholder `json:"dyamicObjs,omitempty"`
	FriendMids []string    `json:"friendMids,omitempty"`
	TitleText  string      `json:"titleText,omitempty"`
	MainText   string      `json:"mainText,omitempty"`
	SubText    string      `json:"subText,omitempty"`
}

type TimelineImage struct {
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type TimelineUrl struct {
	Device    string `json:"device"`
	TargetUrl string `json:"targetUrl"`
}

type TimelineContent struct {
	ApiVer    int               `json:"apiVer"`
	Cmd       string            `json:"cmd"`
	UserMid   string            `json:"userMid"`
	Device    string            `json:"device"`
	Region    string            `json:"region"`
	ChannelID int               `json:"channelID"`
	FeedNo    int               `json:"feedNo"`
	Test      bool              `json:"test,omitempty"`
	PostText  string            `json:"postText"`
	Template  *TimelineTemplate `json:"template,omitempty"`
	Thumbnail *TimelineImage    `json:"thumbnail,omitempty"`
	Url       []*TimelineUrl    `json:"url,omitempty"`
}

type TimelineMsg struct {
	To        []string         `json:"to"`
	ToChannel int              `json:"toChannel"`
	EventType string           `json:"eventType"`
	Content   *TimelineContent `json:"content"`
}
