//go run ./examples/api/main.go -debug -token BOT_TOKEN
package main

import (
	"log"
	"flag"

	"github.com/bot-api/telegram"
	"golang.org/x/net/context"
)


func main() {
	token := flag.String("token", "", "telegram bot token")
	debug := flag.Bool("debug", false, "show debug information")
	flag.Parse()

	if *token == "" {
		log.Fatal("token flag required")
	}

	api := telegram.New(*token)
	api.Debug(*debug)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if user, err := api.GetMe(ctx); err != nil {
		log.Panic(err)
	} else {
		log.Printf("bot info: %#v", user)
	}

	updatesCh := make(chan telegram.Update)

	go telegram.GetUpdates(ctx, api, telegram.UpdateCfg{
		Timeout: 10, 	// Timeout in seconds for long polling.
		Offset: 0, 	// Start with the oldest update
	}, updatesCh)

	for update := range updatesCh {
		log.Printf("got update from %s", update.Message.From.Username)
		if update.Message == nil {
			continue
		}
		msg := telegram.CloneMessage(update.Message, nil)
		// echo with the same message
		if _, err := api.Send(ctx, msg); err != nil {
			log.Printf("send error: %v", err)
		}
	}
}


//go run ./examples/echo/main.go -debug -token BOT_TOKEN
package main
// Simple echo bot, that responses with the same message

import (
	"flag"
	"log"

	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
	"golang.org/x/net/context"
)

func main() {
	token := flag.String("token", "", "telegram bot token")
	debug := flag.Bool("debug", false, "show debug information")
	flag.Parse()

	if *token == "" {
		log.Fatal("token flag is required")
	}

	api := telegram.New(*token)
	api.Debug(*debug)
	bot := telebot.NewWithAPI(api)
	bot.Use(telebot.Recover()) // recover if handler panic

	netCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot.HandleFunc(func(ctx context.Context) error {
		update := telebot.GetUpdate(ctx) // take update from context
		if update.Message == nil {
			return nil
		}
		api := telebot.GetAPI(ctx) // take api from context
		msg := telegram.CloneMessage(update.Message, nil)
		_, err := api.Send(ctx, msg)
		return err

	})

	// Use command middleware, that helps to work with commands
	bot.Use(telebot.Commands(map[string]telebot.Commander{
		"start": telebot.CommandFunc(
			func(ctx context.Context, arg string) error {

				api := telebot.GetAPI(ctx)
				update := telebot.GetUpdate(ctx)
				_, err := api.SendMessage(ctx,
					telegram.NewMessagef(update.Chat().ID,
						"received start with arg %s", arg,
					))
				return err
			}),
	}))


	err := bot.Serve(netCtx)
	if err != nil {
		log.Fatal(err)
	}
}


//go run ./examples/callback/main.go -debug -token BOT_TOKEN

bot.HandleFunc(func(ctx context.Context) error {
    update := telebot.GetUpdate(ctx) // take update from context
    api := telebot.GetAPI(ctx) // take api from context

    if update.CallbackQuery != nil {
        data := update.CallbackQuery.Data
        if strings.HasPrefix(data, "sex:") {
            cfg := telegram.NewEditMessageText(
                update.Chat().ID,
                update.CallbackQuery.Message.MessageID,
                fmt.Sprintf("You sex: %s", data[4:]),
            )
            api.AnswerCallbackQuery(
                ctx,
                telegram.NewAnswerCallback(
                    update.CallbackQuery.ID,
                    "Your configs changed",
                ),
            )
            _, err := api.EditMessageText(ctx, cfg)
            return err
        }
    }

    msg := telegram.NewMessage(update.Chat().ID,
        "Your sex:")
    msg.ReplyMarkup = telegram.InlineKeyboardMarkup{
        InlineKeyboard: telegram.NewVInlineKeyboard(
            "sex:",
            []string{"Female", "Male",},
            []string{"female", "male",},
        ),
    }
    _, err := api.SendMessage(ctx, msg)
    return err

})




//View Source
const (
	// APIEndpoint is the endpoint for all API methods,
	// with formatting for Sprintf.
	APIEndpoint = "https://api.telegram.org/bot%s/%s"
	// FileEndpoint is the endpoint for downloading a file from Telegram.
	FileEndpoint = "https://api.telegram.org/file/bot%s/%s"
)

//View Source
const (
	PrivateChatType    = "private"
	GroupChatType      = "group"
	SuperGroupChatType = "supergroup"
	ChannelChatType    = "channel"
)

//View Source
const (
	ActionTyping         = "typing"
	ActionUploadPhoto    = "upload_photo"
	ActionRecordVideo    = "record_video"
	ActionUploadVideo    = "upload_video"
	ActionRecordAudio    = "record_audio"
	ActionUploadAudio    = "upload_audio"
	ActionUploadDocument = "upload_document"
	ActionFindLocation   = "find_location"
)

//Choose one, depending on what the user is about to receive:
//
typing for text messages
upload_photo for photos
record_video or upload_video for videos
record_audio or upload_audio for audio files
upload_document for general files
find_location for location data


//View Source
const (
	MarkdownMode = "Markdown"
	HTMLMode     = "HTML"
)

//View Source
const (
	// @username
	MentionEntityType     = "mention"
	HashTagEntityType     = "hashtag"
	BotCommandEntityType  = "bot_command"
	URLEntityType         = "url"
	EmailEntityType       = "email"
	BoldEntityType        = "bold"         // bold text
	ItalicEntityType      = "italic"       // italic text
	CodeEntityType        = "code"         // monowidth string
	PreEntityType         = "pre"          // monowidth block
	TextLinkEntityType    = "text_link"    // for clickable text URLs
	TextMentionEntityType = "text_mention" // for users without usernames
)

//View Source
const (
	MemberStatus              = "member"
	CreatorMemberStatus       = "creator"
	AdministratorMemberStatus = "administrator"
	LeftMemberStatus          = "left"
	KickedMemberStatus        = "kicked"
)

//View Source
var DefaultDebugFunc = func(msg string, fields map[string]interface{}) {
	log.Printf("%s %v", msg, fields)
}

//
//      
func GetUpdates(
	ctx context.Context,
	api *API,
	cfg UpdateCfg,
	out chan<- Update) error

//
//      
func IsAPIError(err error) bool

//
//      
func IsForbiddenError(err error) bool

//
//      
func IsRequiredError(err error) bool

//
//      
func IsUnauthorizedError(err error) bool

//
//      
func IsValidToken(token string) bool

//
//      
func IsValidationError(err error) bool

//
//      
func NewHInlineKeyboard(prefix string, text []string, data []string) [][]InlineKeyboardButton

//
//      
func NewHKeyboard(buttons ...string) [][]KeyboardButton

//
//      
func NewKeyboard(buttons [][]string) [][]KeyboardButton

//
//      
func NewVInlineKeyboard(prefix string, text []string, data []string) [][]InlineKeyboardButton

//
//      
func NewVKeyboard(buttons ...string) [][]KeyboardButton

//
//      
type API struct {
	// contains filtered or unexported fields
}

//
//      
func New(token string) *API

//
//      
func NewWithClient(token string, client HTTPDoer) *API

//
//      
func (c *API) AnswerCallbackQuery(
	ctx context.Context,
	cfg AnswerCallbackCfg) (bool, error)

//
//      
func (c *API) AnswerInlineQuery(ctx context.Context, cfg AnswerInlineQueryCfg) (bool, error)

//
//      
func (c *API) Debug(val bool)

//
//      
func (c *API) DebugFunc(f DebugFunc)

//
//      
func (c *API) DownloadFile(ctx context.Context, cfg FileCfg, w io.Writer) error

//
//      
func (c *API) Edit(ctx context.Context, cfg Method) (*EditResult, error)

//
//      
func (c *API) EditMessageCaption(
	ctx context.Context,
	cfg EditMessageCaptionCfg) (*EditResult, error)

//
//      
func (c *API) EditMessageReplyMarkup(
	ctx context.Context,
	cfg EditMessageReplyMarkupCfg) (*EditResult, error)

//
//      
func (c *API) EditMessageText(
	ctx context.Context,
	cfg EditMessageTextCfg) (*EditResult, error)

//
//      
func (c *API) ForwardMessage(
	ctx context.Context,
	cfg ForwardMessageCfg) (*Message, error)

//
//      
func (c *API) GetChat(ctx context.Context, cfg GetChatCfg) (*Chat, error)

//
//      
func (c *API) GetChatAdministrators(
	ctx context.Context,
	cfg GetChatAdministratorsCfg) ([]ChatMember, error)

//
//      
func (c *API) GetChatMember(
	ctx context.Context,
	cfg GetChatMemberCfg) (*ChatMember, error)

//
//      
func (c *API) GetChatMembersCount(
	ctx context.Context,
	cfg GetChatMembersCountCfg) (int, error)

//
//      
func (c *API) GetFile(ctx context.Context, cfg FileCfg) (*File, error)

//
//      
func (c *API) GetMe(ctx context.Context) (*User, error)

//
//      
func (c *API) GetUpdates(
	ctx context.Context,
	cfg UpdateCfg) ([]Update, error)

//
//      
func (c *API) GetUserProfilePhotos(
	ctx context.Context,
	cfg UserProfilePhotosCfg) (*UserProfilePhotos, error)

//
//      
func (c *API) Invoke(ctx context.Context, m Method, dst interface{}) error

//
//      
func (c *API) KickChatMember(
	ctx context.Context,
	cfg KickChatMemberCfg) (bool, error)

//
//      
func (c *API) LeaveChat(
	ctx context.Context,
	cfg LeaveChatCfg) (bool, error)

//
//      
func (c *API) Send(ctx context.Context, cfg Messenger) (*Message, error)

//
//      
func (c *API) SendAudio(
	ctx context.Context,
	cfg AudioCfg) (*Message, error)

//
//      
func (c *API) SendChatAction(ctx context.Context, cfg ChatActionCfg) error

//
//      
func (c *API) SendContact(
	ctx context.Context,
	cfg ContactCfg) (*Message, error)

//
//      
func (c *API) SendDocument(
	ctx context.Context,
	cfg DocumentCfg) (*Message, error)

//
//      
func (c *API) SendMessage(
	ctx context.Context,
	cfg MessageCfg) (*Message, error)

//
//      
func (c *API) SendPhoto(
	ctx context.Context,
	cfg PhotoCfg) (*Message, error)

//
//      
func (c *API) SendSticker(
	ctx context.Context,
	cfg StickerCfg) (*Message, error)

//
//      
func (c *API) SendVenue(
	ctx context.Context,
	cfg VenueCfg) (*Message, error)

//
//      
func (c *API) SendVideo(
	ctx context.Context,
	cfg VideoCfg) (*Message, error)

//
//      
func (c *API) SendVoice(
	ctx context.Context,
	cfg VoiceCfg) (*Message, error)

//
//      
func (c *API) SetWebhook(ctx context.Context, cfg WebhookCfg) error

//
//      
func (c *API) UnbanChatMember(
	ctx context.Context,
	cfg UnbanChatMemberCfg) (bool, error)

//
//      
type APIError struct {
	Description string `json:"description"`
	// ErrorCode contents are subject to change in the future.
	ErrorCode int `json:"error_code"`
}

//
//      
func (e *APIError) Error() string

//
//      
type APIResponse struct {
	Ok          bool             `json:"ok"`
	Result      *json.RawMessage `json:"result"`
	ErrorCode   int              `json:"error_code,omitempty"`
	Description string           `json:"description,omitempty"`
}

//
//      
type AnswerCallbackCfg struct {
	CallbackQueryID string `json:"callback_query_id"`
	Text            string `json:"text"`
	ShowAlert       bool   `json:"show_alert"`
}

//
//      
func NewAnswerCallback(id, text string) AnswerCallbackCfg

//
//      
func NewAnswerCallbackWithAlert(id, text string) AnswerCallbackCfg

//
//      
func (cfg AnswerCallbackCfg) Name() string

//
//      
func (cfg AnswerCallbackCfg) Values() (url.Values, error)

//
//      
type AnswerInlineQueryCfg struct {
	// Unique identifier for the answered query
	InlineQueryID string              `json:"inline_query_id"`
	Results       []InlineQueryResult `json:"results"`
	// The maximum amount of time in seconds
	// that the result of the inline query may be cached on the server.
	// Defaults to 300.
	CacheTime int `json:"cache_time,omitempty"`
	// Pass True, if results may be cached on the server side
	// only for the user that sent the query.
	// By default, results may be returned to any user
	// who sends the same query
	IsPersonal bool `json:"is_personal,omitempty"`
	// Pass the offset that a client should send in the next query
	// with the same text to receive more results.
	// Pass an empty string if there are no more results
	// or if you don‘t support pagination.
	// Offset length can’t exceed 64 bytes.
	NextOffset string `json:"next_offset,omitempty"`
	// If passed, clients will display a button with specified text
	// that switches the user to a private chat with the bot and
	// sends the bot a start message with the parameter switch_pm_parameter
	SwitchPMText string `json:"switch_pm_text,omitempty"`
	// Parameter for the start message sent to the bot
	// when user presses the switch button
	SwitchPMParameter string `json:"switch_pm_parameter"`
}

//
//      
func (cfg AnswerInlineQueryCfg) Name() string

//
//      
func (cfg AnswerInlineQueryCfg) Values() (url.Values, error)

//
//      
type Audio struct {
	MetaFile

	// Duration of the recording in seconds as defined by sender.
	Duration int `json:"duration"`
	//  Performer of the audio as defined by sender
	// or by audio tags. Optional.
	Performer string `json:"performer,omitempty"`
	// Title of the audio as defined by sender
	// or by audio tags. Optional.
	Title string `json:"title,omitempty"`
	// MIMEType of the file as defined by sender. Optional.
	MIMEType string `json:"mime_type,omitempty"`
}

//
//      
type AudioCfg struct {
	BaseFile
	Duration  int
	Performer string
	Title     string
}

//
//      
func (cfg AudioCfg) Field() string

//
//      
func (cfg AudioCfg) Name() string

//
//      
func (cfg AudioCfg) Values() (url.Values, error)

//
//      
type BaseChat struct {
	// Unique identifier for the target chat
	ID int64
	// Username of the target channel (in the format @channelusername)
	ChannelUsername string
}

//
//      
func (c *BaseChat) SetChatID(id int64)

//
//      
func (c BaseChat) Values() (url.Values, error)

//
//      
type BaseEdit struct {
	// Required if inline_message_id is not specified.
	// Unique identifier for the target chat or
	// username of the target channel (in the format @channelusername)
	ChatID          int64
	ChannelUsername string
	// Required if inline_message_id is not specified.
	// Unique identifier of the sent message
	MessageID int64
	// Required if chat_id and message_id are not specified.
	// Identifier of the inline message
	InlineMessageID string
	// Only InlineKeyboardMarkup supported right now.
	ReplyMarkup ReplyMarkup
}

//
//      
func (m BaseEdit) Values() (url.Values, error)

//
//      
type BaseFile struct {
	BaseMessage
	FileID    string
	MimeType  string
	InputFile InputFile
}

//
//      
func (b BaseFile) Exist() bool

//
//      
func (b BaseFile) File() InputFile

//
//      
func (b BaseFile) GetFileID() string

//
//      
func (b *BaseFile) Reset(i InputFile)

//
//      
func (b BaseFile) Values() (url.Values, error)

//
//      
type BaseInlineQueryResult struct {
	Type                string              `json:"type"` // required
	ID                  string              `json:"id"`   // required
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
	// ReplyMarkup supports only InlineKeyboardMarkup for InlineQueryResult
	ReplyMarkup ReplyMarkup `json:"reply_markup,omitempty"`
}

//
//      
type BaseMessage struct {
	BaseChat
	// If the message is a reply, ID of the original message
	ReplyToMessageID int64
	// Additional interface options.
	// A JSON-serialized object for a custom reply keyboard,
	// instructions to hide keyboard or to force a reply from the user.
	ReplyMarkup ReplyMarkup
	// Sends the message silently.
	// iOS users will not receive a notification,
	// Android users will receive a notification with no sound.
	// Other apps coming soon.
	DisableNotification bool
}

//
//      
func (BaseMessage) Message() *Message

//
//      
func (m BaseMessage) Values() (url.Values, error)

//
//      
type CallbackQuery struct {
	// Unique identifier for this query
	ID string `json:"id"`
	// Sender
	From *User `json:"from"`
	// Message with the callback button that originated the query.
	// Note that message content and message date
	// will not be available if the message is too old. Optional.
	Message *Message `json:"message,omitempty"`
	// Identifier of the message sent via the bot in inline mode,
	// that originated the query. Optional.
	InlineMessageID string `json:"inline_message_id,omitempty"`
	// Data associated with the callback button.
	// Be aware that a bad client can send arbitrary data in this field
	Data string `json:"data"`
}

//
//      
type Chat struct {
	// ID is a Unique identifier for this chat, not exceeding 1e13 by absolute value.
	ID int64 `json:"id"`
	// Type of chat, can be either “private”, “group”, "supergroup" or “channel”
	Type string `json:"type"`

	Title     string `json:"title,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
}

//
//      
type ChatActionCfg struct {
	BaseChat
	// Type of action to broadcast.
	// Choose one, depending on what the user is about to receive:
	// typing for text messages, upload_photo for photos,
	// record_video or upload_video for videos,
	// record_audio or upload_audio for audio files,
	// upload_document for general files,
	// find_location for location data.
	// Use one of constants: ActionTyping, ActionFindLocation, etc
	Action string
}

//
//      
func NewChatAction(chatID int64, action string) ChatActionCfg

//
//      
func (cfg ChatActionCfg) Name() string

//
//      
func (cfg ChatActionCfg) Values() (url.Values, error)

//
//      
type ChatMember struct {
	// Information about the user.
	User User `json:"user"`
	// The member's status in the chat.
	// One of MemberStatus constants.
	Status string `json:"status"`
}

//
//      
type ChosenInlineResult struct {
	// ResultID is a unique identifier for the result that was chosen.
	ResultID string `json:"result_id"`
	// From is a user that chose the result.
	From User `json:"from"`
	// Query is used to obtain the result.
	Query string `json:"query"`
}

//
//      
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`

	// UserID is a contact's user identifier in Telegram. Optional.
	UserID   int64  `json:"user_id,omitempty"`
	LastName string `json:"last_name,omitempty"`
}

//
//      
func (c Contact) Values() url.Values

//
//      
type ContactCfg struct {
	BaseMessage
	Contact
}

//
//      
func (cfg ContactCfg) Name() string

//
//      
func (cfg ContactCfg) Values() (url.Values, error)

//
//      
type DebugFunc func(msg string, fields map[string]interface{})

//
//      
type Document struct {
	MetaFile

	// Document thumbnail as defined by sender. Optional.
	Thumb *PhotoSize `json:"thumb,omitempty"`

	// Original filename as defined by sender. Optional.
	FileName string `json:"file_name,omitempty"`

	// MIMEType of the file as defined by sender. Optional.
	MIMEType string `json:"mime_type,omitempty"`
}

//
//      
type DocumentCfg struct {
	BaseFile
}

//
//      
func (cfg DocumentCfg) Field() string

//
//      
func (cfg DocumentCfg) Name() string

//
//      
func (cfg DocumentCfg) Values() (url.Values, error)

//
//      
type EditMessageCaptionCfg struct {
	BaseEdit
	// New caption of the message
	Caption string
}

//
//      
func NewEditMessageCaption(chatID, messageID int64, caption string) EditMessageCaptionCfg

//
//      
func (EditMessageCaptionCfg) Name() string

//
//      
func (cfg EditMessageCaptionCfg) Values() (url.Values, error)

//
//      
type EditMessageReplyMarkupCfg struct {
	BaseEdit
}

//
//      
func NewEditMessageReplyMarkup(chatID, messageID int64, replyMarkup *InlineKeyboardMarkup) EditMessageReplyMarkupCfg

//
//      
func (EditMessageReplyMarkupCfg) Name() string

//
//      
func (cfg EditMessageReplyMarkupCfg) Values() (url.Values, error)

//
//      
type EditMessageTextCfg struct {
	BaseEdit
	// New text of the message
	Text string
	// Send Markdown or HTML, if you want Telegram apps
	// to show bold, italic, fixed-width text
	// or inline URLs in your bot's message. Optional.
	ParseMode string
	// Disables link previews for links in this message. Optional.
	DisableWebPagePreview bool
}

//
//      
func NewEditMessageText(chatID, messageID int64, text string) EditMessageTextCfg

//
//      
func (EditMessageTextCfg) Name() string

//
//      
func (cfg EditMessageTextCfg) Values() (url.Values, error)

//
//      
type EditResult struct {
	Message *Message
	Ok      bool
}

//
//      
func (e *EditResult) UnmarshalJSON(data []byte) error

//
//      
type File struct {
	MetaFile
	// FilePath is a relative path to file.
	// Use https://api.telegram.org/file/bot<token>/<file_path>
	// to get the file.
	FilePath string `json:"file_path,omitempty"`

	// Link is inserted by Api client after GetFile request
	Link string `json:"link"`
}

//
//      
type FileCfg struct {
	FileID string
}

//
//      
func (cfg FileCfg) Name() string

//
//      
func (cfg FileCfg) Values() (url.Values, error)

//
//      
type Filer interface {
	// Field name for file data
	Field() string
	// File data
	File() InputFile
	// Exist returns true if file exists on telegram servers
	Exist() bool
	// Reset removes FileID and sets new InputFile
	// Reset(InputFile)
	// GetFileID returns fileID if it's exist
	GetFileID() string
}

//
//      
type ForceReply struct {
	MarkReplyMarkup

	ForceReply bool `json:"force_reply"`
	Selective  bool `json:"selective"` // optional
}

//
//      
type ForwardMessageCfg struct {
	BaseChat
	// Unique identifier for the chat where the original message was sent
	FromChat BaseChat
	// Unique message identifier
	MessageID int64
	// Sends the message silently.
	// iOS users will not receive a notification,
	// Android users will receive a notification with no sound.
	// Other apps coming soon.
	DisableNotification bool
}

//
//      
func NewForwardMessage(chatID, fromChatID, messageID int64) ForwardMessageCfg

//
//      
func (ForwardMessageCfg) Message() *Message

//
//      
func (cfg ForwardMessageCfg) Name() string

//
//      
func (cfg ForwardMessageCfg) Values() (url.Values, error)

//
//      
type GetChatAdministratorsCfg struct {
	BaseChat
}

//
//      
func (cfg GetChatAdministratorsCfg) Name() string

//
//      
func (cfg GetChatAdministratorsCfg) Values() (url.Values, error)

//
//      
type GetChatCfg struct {
	BaseChat
}

//
//      
func (cfg GetChatCfg) Name() string

//
//      
func (cfg GetChatCfg) Values() (url.Values, error)

//
//      
type GetChatMemberCfg struct {
	BaseChat
	UserID int64 `json:"user_id"`
}

//
//      
func (cfg GetChatMemberCfg) Name() string

//
//      
func (cfg GetChatMemberCfg) Values() (url.Values, error)

//
//      
type GetChatMembersCountCfg struct {
	BaseChat
}

//
//      
func (cfg GetChatMembersCountCfg) Name() string

//
//      
func (cfg GetChatMembersCountCfg) Values() (url.Values, error)

//
//      
type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

//
//      
type InlineKeyboardButton struct {
	// Label text on the button
	Text string `json:"text"`
	// HTTP url to be opened when button is pressed. Optional.
	URL string `json:"url,omitempty"`
	// Data to be sent in a callback query to the bot
	// when button is pressed. Optional.
	CallbackData string `json:"callback_data,omitempty"`
	// If set, pressing the button will prompt the user
	// to select one of their chats, open that chat
	// and insert the bot‘s username and the specified inline query
	// in the input field. Can be empty,
	// in which case just the bot’s username will be inserted.
	// Optional.
	//
	// Note: This offers an easy way for users
	// to start using your bot in inline mode
	// when they are currently in a private chat with it.
	// Especially useful when combined with switch_pm... actions
	// – in this case the user will be automatically returned to the chat
	// they switched from, skipping the chat selection screen.
	SwitchInlineQuery string `json:"switch_inline_query,omitempty"`
}

//
//      
type InlineKeyboardMarkup struct {
	MarkReplyMarkup

	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

//
//      
type InlineQuery struct {
	// ID is a unique identifier for this query.
	ID string `json:"id"`
	// From is a sender.
	From User `json:"from"`
	// Sender location, only for bots that request user location.
	// Optional.
	Location *Location `json:"location,omitempty"`
	// Query is a text of the query.
	Query string `json:"query"`
	// Offset of the results to be returned, can be controlled by the bot.
	Offset string `json:"offset"`
}

//
//      
type InlineQueryResult interface {
	// InlineQueryResult is a fake method that helps to identify implementations
	InlineQueryResult()
}

//
//      
type InlineQueryResultArticle struct {
	MarkInlineQueryResult
	BaseInlineQueryResult
	InlineThumb

	Title string `json:"title"` // required
	URL   string `json:"url,omitempty"`
	// Optional. Pass True, if you don't want the URL
	// to be shown in the message
	HideURL     bool   `json:"hide_url,omitempty"`
	Description string `json:"description,omitempty"`
}

//
//      
func NewInlineQueryResultArticle(id, title, messageText string) *InlineQueryResultArticle

//
//      
type InlineQueryResultAudio struct {
	MarkInlineQueryResult
	BaseInlineQueryResult

	AudioURL      string `json:"audio_url"` // required
	Title         string `json:"title"`     // required
	Performer     string `json:"performer"`
	AudioDuration int    `json:"audio_duration"`
}

//
//      
type InlineQueryResultContact struct {
	MarkInlineQueryResult
	BaseInlineQueryResult
	InlineThumb
	Contact
}

//
//      
type InlineQueryResultDocument struct {
	MarkInlineQueryResult
	BaseInlineQueryResult
	InlineThumb

	DocumentURL string `json:"document_url"` // required
	// Mime type of the content of the file,
	// either “application/pdf” or “application/zip”
	MimeType string `json:"mime_type"` // required
	// Title for the result
	Title string `json:"title"` // required
	// Optional. Caption of the document to be sent, 0-200 characters
	Caption string `json:"caption,omitempty"`
	// Optional. Short description of the result
	Description string `json:"description,omitempty"`
}

//
//      
type InlineQueryResultGIF struct {
	MarkInlineQueryResult
	BaseInlineQueryResult
	InlineThumb

	// A valid URL for the GIF file. File size must not exceed 1MB
	GifURL    string `json:"gif_url"` // required
	GifWidth  int    `json:"gif_width,omitempty"`
	GifHeight int    `json:"gif_height,omitempty"`
	Title     string `json:"title,omitempty"`
	Caption   string `json:"caption,omitempty"`
}

//
//      
type InlineQueryResultLocation struct {
	MarkInlineQueryResult
	BaseInlineQueryResult
	InlineThumb

	Latitude  float64 `json:"latitude"`  // required
	Longitude float64 `json:"longitude"` // required
	Title     string  `json:"title"`     // required
}

//
//      
type InlineQueryResultMPEG4GIF struct {
	MarkInlineQueryResult
	BaseInlineQueryResult

	MPEG4URL    string `json:"mpeg4_url"` // required
	MPEG4Width  int    `json:"mpeg4_width,omitempty"`
	MPEG4Height int    `json:"mpeg4_height,omitempty"`
	Title       string `json:"title,omitempty"`
	Caption     string `json:"caption,omitempty"`
}

//
//      
type InlineQueryResultPhoto struct {
	MarkInlineQueryResult
	BaseInlineQueryResult
	InlineThumb

	PhotoURL    string `json:"photo_url"` // required
	PhotoWidth  int    `json:"photo_width,omitempty"`
	PhotoHeight int    `json:"photo_height,omitempty"`
	MimeType    string `json:"mime_type,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Caption     string `json:"caption,omitempty"`
}

//
//      
type InlineQueryResultVenue struct {
	MarkInlineQueryResult
	BaseInlineQueryResult
	InlineThumb
	Venue
}

//
//      
type InlineQueryResultVideo struct {
	MarkInlineQueryResult
	BaseInlineQueryResult
	InlineThumb

	VideoURL      string `json:"video_url"` // required
	MimeType      string `json:"mime_type"` // required
	Title         string `json:"title,omitempty"`
	Caption       string `json:"caption,omitempty"`
	VideoWidth    int    `json:"video_width,omitempty"`
	VideoHeight   int    `json:"video_height,omitempty"`
	VideoDuration int    `json:"video_duration,omitempty"`
	Description   string `json:"description,omitempty"`
}

//
//      
type InlineQueryResultVoice struct {
	MarkInlineQueryResult
	BaseInlineQueryResult

	VoiceURL string `json:"voice_url"` // required
	Title    string `json:"title"`     // required
	Duration int    `json:"voice_duration,omitempty"`
}

//
//      
type InlineThumb struct {
	ThumbURL    string `json:"thumb_url,omitempty"`
	ThumbWidth  int    `json:"thumb_width,omitempty"`
	ThumbHeight int    `json:"thumb_height,omitempty"`
}

//
//      
type InputContactMessageContent struct {
	MarkInputMessageContent
	Contact
}

//
//      
type InputFile interface {
	Reader() io.Reader
	Name() string
}

//
//      
func NewBytesFile(filename string, data []byte) InputFile

//
//      
func NewInputFile(filename string, r io.Reader) InputFile

//
//      
type InputLocationMessageContent struct {
	MarkInputMessageContent
	Location
}

//
//      
type InputMessageContent interface {
	// MessageContent is a fake method that helps to identify implementations
	InputMessageContent()
}

//
//      
type InputTextMessageContent struct {
	MarkInputMessageContent

	// Text of the message to be sent, 1‐4096 characters
	MessageText string `json:"message_text"`
	// Send Markdown or HTML, if you want Telegram apps to show
	// bold, italic, fixed‐width text or inline URLs in your bot's message.
	// Use Mode constants. Optional.
	ParseMode string `json:"parse_mode,omitempty"`
	// Disables link previews for links in this message.
	DisableWebPagePreview bool `json:"disable_web_page_preview,omitempty"`
}

//
//      
type InputVenueMessageContent struct {
	MarkInputMessageContent
	Venue
}

//
//      
type KeyboardButton struct {
	// Text of the button. If none of the optional fields are used,
	// it will be sent to the bot as a message when the button is pressed
	Text string `json:"text"`
	// If true, the user's phone number will be sent as a contact
	// when the button is pressed. Available in private chats only.
	// Optional.
	RequestContact bool `json:"request_contact,omitempty"`
	// If true, the user's current location will be sent
	// when the button is pressed. Available in private chats only.
	// Optional.
	RequestLocation bool `json:"request_location,omitempty"`
}

//
//      
type KickChatMemberCfg struct {
	BaseChat
	UserID int64 `json:"user_id"`
}

//
//      
func (cfg KickChatMemberCfg) Name() string

//
//      
func (cfg KickChatMemberCfg) Values() (url.Values, error)

//
//      
type LeaveChatCfg struct {
	BaseChat
}

//
//      
func (cfg LeaveChatCfg) Name() string

//
//      
func (cfg LeaveChatCfg) Values() (url.Values, error)

//
//      
type Location struct {
	// Longitude as defined by sender
	Longitude float64 `json:"longitude"`
	// Latitude as defined by sender
	Latitude float64 `json:"latitude"`
}

//
//      
func (l Location) Values() url.Values

//
//      
type LocationCfg struct {
	BaseMessage
	Location
}

//
//      
func NewLocation(chatID int64, lat float64, lon float64) LocationCfg

//
//      
func (cfg LocationCfg) Name() string

//
//      
func (cfg LocationCfg) Values() (url.Values, error)

//
//      
type MarkInlineQueryResult struct{}

//
//      
func (MarkInlineQueryResult) InlineQueryResult()

//
//      
type MarkInputMessageContent struct{}

//
//      
func (MarkInputMessageContent) InputMessageContent()

//
//      
type MarkReplyMarkup struct{}

//
//      
func (MarkReplyMarkup) ReplyMarkup()

//
//      
type MeCfg struct{}

//
//      
func (cfg MeCfg) Name() string

//
//      
func (cfg MeCfg) Values() (url.Values, error)

//
//      
type Message struct {
	// MessageID is a unique message identifier.
	MessageID int64 `json:"message_id"`

	// From is a sender, can be empty for messages sent to channels.
	// Optional.
	From *User `json:"from,omitempty"`
	// Date  the message was sent in Unix time.
	Date int `json:"date"`
	// Chat is a conversation the message belongs to.
	Chat Chat `json:"chat"`

	// ForwardFrom is a sender of the original message
	// for forwarded messages. Optional.
	ForwardFrom *User `json:"forward_from,omitempty"`

	// For messages forwarded from a channel,
	// information about the original channel. Optional.
	ForwardFromChat *Chat `json:"forward_from_chat,omitempty"`

	// ForwardDate is a unixtime of the original message
	// for forwarded messages. Optional.
	ForwardDate int `json:"forward_date"`

	// ReplyToMessage is an original message for replies.
	// Note that the Message object in this field will not
	// contain further ReplyToMessage fields even if it
	// itself is a reply. Optional.
	ReplyToMessage *Message `json:"reply_to_message,omitempty"`

	// Date the message was last edited in Unix time
	// Zero value means object wasn't edited.
	// Optional.
	EditDate int `json:"edit_date,omitempty"`

	// For text messages, special entities like usernames,
	// URLs, bot commands, etc. that appear in the text. Optional
	Entities []MessageEntity `json:"entities"`

	// Text is an actual UTF-8 text of the message for a text message,
	// 0-4096 characters. Optional.
	Text string `json:"text,omitempty"`

	// Audio has information about the audio file. Optional.
	Audio *Audio `json:"audio,omitempty"`

	// Document has information about a general file. Optional.
	Document *Document `json:"document,omitempty"`

	// Photo has a slice of available sizes of photo. Optional.
	Photo []PhotoSize `json:"photo,omitempty"`

	// Sticker has information about the sticker. Optional.
	Sticker *Sticker `json:"sticker,omitempty"`

	// For a video, information about it.
	Video *Video `json:"video,omitempty"`

	// Message is a voice message, information about the file
	Voice *Voice `json:"voice,omitempty"`

	// Caption for the document, photo or video, 0‐200 characters
	Caption string `json:"caption,omitempty"`

	// For a contact, contact information itself.
	Contact *Contact `json:"contact,omitempty"`

	// For a location, its longitude and latitude.
	Location *Location `json:"location,omitempty"`

	// Message is a venue, information about the venue
	Venue *Venue `json:"venue,omitempty"`

	// NewChatMember has an information about a new member
	// that was added to the group
	// (this member may be the bot itself). Optional.
	NewChatMember *User `json:"new_chat_member,omitempty"`

	// LeftChatMember has an information about a member
	// that was removed from the group
	// (this member may be the bot itself). Optional.
	LeftChatMember *User `json:"left_chat_member,omitempty"`

	// For a service message, represents a new title
	// for chat this message came from.
	//
	// Sender would lead to a User, capable of change.
	NewChatTitle string `json:"new_chat_title,omitempty"`

	// For a service message, represents all available
	// thumbnails of new chat photo.
	//
	// Sender would lead to a User, capable of change.
	NewChatPhoto []PhotoSize `json:"new_chat_photo,omitempty"`

	// For a service message, true if chat photo just
	// got removed.
	//
	// Sender would lead to a User, capable of change.
	DeleteChatPhoto bool `json:"delete_chat_photo,omitempty"`

	// For a service message, true if group has been created.
	//
	// You would receive such a message if you are one of
	// initial group chat members.
	//
	// Sender would lead to creator of the chat.
	GroupChatCreated bool `json:"group_chat_created,omitempty"`

	// For a service message, true if super group has been created.
	//
	// You would receive such a message if you are one of
	// initial group chat members.
	//
	// Sender would lead to creator of the chat.
	SuperGroupChatCreated bool `json:"supergroup_chat_created,omitempty"`

	// For a service message, true if channel has been created.
	//
	// You would receive such a message if you are one of
	// initial channel administrators.
	//
	// Sender would lead to creator of the chat.
	ChannelChatCreated bool `json:"channel_chat_created,omitempty"`

	// For a service message, the destination (super group) you
	// migrated to.
	//
	// You would receive such a message when your chat has migrated
	// to a super group.
	//
	// Sender would lead to creator of the migration.
	MigrateToChatID int64 `json:"migrate_to_chat_id,omitempty"`

	// For a service message, the Origin (normal group) you migrated
	// from.
	//
	// You would receive such a message when your chat has migrated
	// to a super group.
	//
	// Sender would lead to creator of the migration.
	MigrateFromChatID int64 `json:"migrate_from_chat_id,omitempty"`
	// Specified message was pinned.
	// Note that the Message object in this field
	// will not contain further reply_to_message fields
	// even if it is itself a reply.
	PinnedMessage *Message `json:"pinned_message,omitempty"`
}

//
//      
func (m *Message) Command() (string, string)

//
//      
func (m *Message) IsCommand() bool

//
//      
type MessageCfg struct {
	BaseMessage
	Text string
	// Send Markdown or HTML, if you want Telegram apps to show
	// bold, italic, fixed-width text or inline URLs in your bot's message.
	// Use one of constants: ModeHTML, ModeMarkdown.
	ParseMode string
	// Disables link previews for links in this message.
	DisableWebPagePreview bool
}

//
//      
func NewMessage(chatID int64, text string) MessageCfg

//
//      
func NewMessagef(chatID int64, text string, args ...interface{}) MessageCfg

//
//      
func (cfg MessageCfg) Name() string

//
//      
func (cfg MessageCfg) Values() (url.Values, error)

//
//      
type MessageEntity struct {
	// Type of the entity. One of mention ( @username ), hashtag,
	// bot_command, url, email, bold (bold text),
	// italic (italic text), code (monowidth string),
	// pre (monowidth block), text_link (for clickable text URLs),
	// text_mention (for users without usernames)
	// Use constants SomethingEntityType instead of string.
	Type string `json:"type"`
	// Offset in UTF‐16 code units to the start of the entity
	Offset int `json:"offset"`
	// Length of the entity in UTF‐16 code units
	Length int `json:"length"`
	// For “text_link” only, url that will be opened
	// after user taps on the text. Optional
	URL string `json:"url,omitempty"`
	// For “text_mention” only, the mentioned user. Optional.
	User *User `json:"user,omitempty"`
}

//
//      
type Messenger interface {
	Method
	Message() *Message
}

//
//      
func CloneMessage(msg *Message, baseMessage *BaseMessage) Messenger

//
//      
type MetaFile struct {
	// FileID is a Unique identifier for this file.
	FileID string `json:"file_id"`
	// FileSize is a size of file if known. Optional.
	FileSize int `json:"file_size,omitempty"`
}

//
//      
type Method interface {
	// method name
	Name() string
	// method params
	Values() (url.Values, error)
}

//
//      
type PhotoCfg struct {
	BaseFile
	Caption string
}

//
//      
func NewPhotoShare(chatID int64, fileID string) PhotoCfg

//
//      
func NewPhotoUpload(chatID int64, inputFile InputFile) PhotoCfg

//
//      
func (cfg PhotoCfg) Field() string

//
//      
func (cfg PhotoCfg) Name() string

//
//      
func (cfg PhotoCfg) Values() (url.Values, error)

//
//      
type PhotoSize struct {
	MetaFile
	Size
}

//
//      
type ReplyKeyboardHide struct {
	MarkReplyMarkup

	HideKeyboard bool `json:"hide_keyboard"`
	Selective    bool `json:"selective"` // optional
}

//
//      
type ReplyKeyboardMarkup struct {
	MarkReplyMarkup

	// Array of button rows, each represented by an Array of Strings
	Keyboard [][]KeyboardButton `json:"keyboard"`
	// Requests clients to resize the keyboard vertically
	// for optimal fit (e.g., make the keyboard smaller
	// if there are just two rows of buttons).
	// Defaults to false, in which case the custom keyboard
	// is always of the same height as the app's standard keyboard.
	ResizeKeyboard bool `json:"resize_keyboard,omitempty"`
	// Requests clients to hide the keyboard as soon as it's been used.
	// The keyboard will still be available,
	// but clients will automatically display the usual
	// letter‐keyboard in the chat – the user can press
	// a special button in the input field to see the custom keyboard again.
	// Defaults to false.
	OneTimeKeyboard bool `json:"one_time_keyboard,omitempty"`
	// Use this parameter if you want to show the keyboard
	// to specific users only.
	// Targets:
	// 1) users that are @mentioned in the text of the Message object;
	// 2) if the bot's message is a reply (has reply_to_message_id),
	// sender of the original message.
	Selective bool `json:"selective,omitempty"`
}

//
//      
type ReplyMarkup interface {
	// ReplyMarkup is a fake method that helps to identify implementations
	ReplyMarkup()
}

//
//      
type RequiredError struct {
	Fields []string
}

//
//      
func NewRequiredError(fields ...string) *RequiredError

//
//      
func (e *RequiredError) Error() string

//
//      
type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

//
//      
type Sticker struct {
	MetaFile
	Size // Sticker width and height

	// Sticker thumbnail in .webp or .jpg format. Optional.
	Thumb *PhotoSize `json:"thumb,omitempty"`
	// Emoji associated with the sticker. Optional.
	Emoji string `json:"emoji,omitempty"`
}

//
//      
type StickerCfg struct {
	BaseFile
}

//
//      
func (cfg StickerCfg) Field() string

//
//      
func (cfg StickerCfg) Name() string

//
//      
func (cfg StickerCfg) Values() (url.Values, error)

//
//      
type UnbanChatMemberCfg struct {
	BaseChat
	UserID int64 `json:"user_id"`
}

//
//      
func (cfg UnbanChatMemberCfg) Name() string

//
//      
func (cfg UnbanChatMemberCfg) Values() (url.Values, error)

//
//      
type Update struct {
	// UpdateID is the update‘s unique identifier.
	// Update identifiers start from a certain positive number
	// and increase sequentially
	UpdateID int64 `json:"update_id"`
	// Message is a new incoming message of any kind:
	// text, photo, sticker, etc. Optional.
	Message *Message `json:"message,omitempty"`
	// New version of a message that is known to the bot and was edited.
	// Optional.
	EditedMessage *Message `json:"edited_message,omitempty"`
	// InlineQuery is a new incoming inline query. Optional.
	InlineQuery *InlineQuery `json:"inline_query,omitempty"`
	// ChosenInlineResult is a result of an inline query
	// that was chosen by a user and sent to their chat partner. Optional.
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result,omitempty"`
	// CallbackQuery is a new incoming callback query. Optional.
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
}

//
//      
func (u Update) Chat() (chat *Chat)

//
//      
func (u Update) From() (from *User)

//
//      
func (u Update) HasMessage() bool

//
//      
func (u Update) IsEdited() bool

//
//      
type UpdateCfg struct {
	// Identifier of the first update to be returned.
	// Must be greater by one than the highest
	// among the identifiers of previously received updates.
	// By default, updates starting with the earliest
	// unconfirmed update are returned. An update is considered confirmed
	// as soon as getUpdates is called with an offset
	// higher than its update_id. The negative offset
	// can be specified to retrieve updates starting
	// from -offset update from the end of the updates queue.
	// All previous updates will forgotten.
	Offset int64
	// Limits the number of updates to be retrieved.
	// Values between 1—100 are accepted. Defaults to 100.
	Limit int
	// Timeout in seconds for long polling.
	// Defaults to 0, i.e. usual short polling
	Timeout int
}

//
//      
func NewUpdate(offset int64) UpdateCfg

//
//      
func (cfg UpdateCfg) Name() string

//
//      
func (cfg UpdateCfg) Values() (url.Values, error)

//
//      
type User struct {
	// ID is a unique identifier for this user or bot.
	ID int64 `json:"id"`
	// FirstName is a user‘s or bot’s first name
	FirstName string `json:"first_name"`
	// LastName is a user‘s or bot’s last name. Optional.
	LastName string `json:"last_name,omitempty"`
	// Username is a user‘s or bot’s username. Optional.
	Username string `json:"username,omitempty"`
}

//
//      
type UserProfilePhotos struct {
	// Total number of profile pictures the target user has
	TotalCount int `json:"total_count"`
	// Requested profile pictures (in up to 4 sizes each)
	Photos [][]PhotoSize `json:"photos"`
}

//
//      
type UserProfilePhotosCfg struct {
	UserID int64
	// Sequential number of the first photo to be returned.
	// By default, all photos are returned.
	Offset int
	// Limits the number of photos to be retrieved.
	// Values between 1—100 are accepted. Defaults to 100.
	Limit int
}

//
//      
func NewUserProfilePhotos(userID int64) UserProfilePhotosCfg

//
//      
func (cfg UserProfilePhotosCfg) Name() string

//
//      
func (cfg UserProfilePhotosCfg) Values() (url.Values, error)

//
//      
type ValidationError struct {
	// Field name
	Field       string `json:"field"`
	Description string `json:"description"`
}

//
//      
func NewValidationError(field string, description string) *ValidationError

//
//      
func (e *ValidationError) Error() string

//
//      
type Venue struct {
	// Venue location
	Location Location `json:"location"`
	// Name of the venue
	Title string `json:"title"`
	// Address of the venue
	Address string `json:"address"`
	// Foursquare identifier of the venue. Optional.
	FoursquareID string `json:"foursquare_id,omitempty"`
}

//
//      
func (venue Venue) Values() url.Values

//
//      
type VenueCfg struct {
	BaseMessage
	Venue
}

//
//      
func (cfg VenueCfg) Name() string

//
//      
func (cfg VenueCfg) Values() (url.Values, error)

//
//      
type Video struct {
	MetaFile
	Size

	// Duration of the recording in seconds as defined by sender.
	Duration int `json:"duration"`
	// MIMEType of the file as defined by sender. Optional.
	MIMEType string `json:"mime_type,omitempty"`
	// Video thumbnail. Optional.
	Thumb *PhotoSize `json:"thumb,omitempty"`
}

//
//      
type VideoCfg struct {
	BaseFile
	Duration int
	Caption  string
}

//
//      
func (cfg VideoCfg) Field() string

//
//      
func (cfg VideoCfg) Name() string

//
//      
func (cfg VideoCfg) Values() (url.Values, error)

//
//      
type Voice struct {
	MetaFile

	// Duration of the recording in seconds as defined by sender.
	Duration int `json:"duration"`
	// MIMEType of the file as defined by sender. Optional.
	MIMEType string `json:"mime_type,omitempty"`
}

//
//      
type VoiceCfg struct {
	BaseFile
	Duration int
}

//
//      
func (cfg VoiceCfg) Field() string

//
//      
func (cfg VoiceCfg) Name() string

//
//      
func (cfg VoiceCfg) Values() (url.Values, error)

//
//      
type WebhookCfg struct {
	URL string
	// self generated TLS certificate
	Certificate InputFile
}

//
//      
func NewWebhook(link string) WebhookCfg

//
//      
func NewWebhookWithCert(link string, file InputFile) WebhookCfg

//
//      
func (cfg WebhookCfg) Exist() bool

//
//      
func (cfg WebhookCfg) Field() string

//
//      
func (cfg WebhookCfg) File() InputFile

//
//      
func (cfg *WebhookCfg) GetFileID() string

//
//      
func (cfg WebhookCfg) Name() string

//
//      
func (cfg *WebhookCfg) Reset(i InputFile)

//
//      
func (cfg WebhookCfg) Values() (url.Values, error)


