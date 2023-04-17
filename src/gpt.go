import (
    "fmt"
    "github.com/bwmarrin/discordgo"
    "github.com/jonas747/dca"
    "github.com/bwmarrin/discordgo/voice"
    "github.com/bwmarrin/discordgo/voice/ffmpeg"
    "github.com/qiangxue/fasthttp-routing"
    "github.com/TeamSpeak-Systems/go-ts3"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

//hold the bot's state
type Bot struct {
    Session *discordgo.Session
    Voice   *voice.Connection
    Ts3     *ts3.Client
}

//handle discord bot events
func (bot *Bot) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == bot.Session.State.User.ID {
        return
    }
 
    if m.Content == "!join" {
        bot.joinVoiceChannel(m.GuildID, m.Author.ID, m.Member.VoiceState.ChannelID)
    }
 
    if m.Content == "!leave" {
        bot.leaveVoiceChannel(m.GuildID)
    }
 
    if m.Content == "!join-ts3" {
        bot.joinTs3Server(m.GuildID, m.Author.ID)
    }
 
    if m.Content == "!leave-ts3" {
        bot.leaveTs3Server(m.GuildID)
    }
 
    if m.Content == "!play" {
        bot.playMusic(m.GuildID, m.Author.ID, m.Content)
    }
}

//handle joining and leaving discord voice channels
func (bot *Bot) joinVoiceChannel(guildID, userID, channelID string) error {
    voiceChannel, err := bot.Session.ChannelVoiceJoin(guildID, channelID, false, true)
    if err != nil {
        return err
    }
 
    bot.Voice = voiceChannel
    return nil
}
func (bot *Bot) leaveVoiceChannel(guildID string) error {
    if bot.Voice == nil {
        return errors.New("Not currently in a voice channel")
    }
 
    err := bot.Voice.Disconnect()
    if err != nil {
        return err
    }
 
    bot.Voice = nil
    return nil
}

//handle joining and leaving teamspeak servers
func (bot *Bot) joinTs3Server(guildID, userID string) error {
    serverAddress := "localhost:9987"
    serverPassword := "password"
    clientName := fmt.Sprintf("%s_%s", guildID, userID)
 
    ts3Client, err := ts3.NewClient(serverAddress)
    if err != nil {
        return err
    }
 
    err = ts3Client.Login(clientName, serverPassword)
    if err != nil {
        return err
    }
 
    err = ts3Client.Use(1)
    if err != nil {
        return err
    }
 
    bot.Ts3 = ts3Client
    return nil
}
func (bot *Bot) leaveTs3Server(guildID string) error {
    if bot.Ts3 == nil {
        return errors.New("Not currently connected to a Teamspeak server")
    }
 
    err := bot.Ts3.Logout()
    if err != nil {
        return err
    }
 
    bot.Ts
  .Close()
bot.Ts3 = nil
return nil
}

//play music from yt and soundcloud links
func (bot *Bot) playMusic(guildID, userID, content string) error {
    if bot.Voice == nil {
        return errors.New("Not currently in a voice channel")
    }
 
    // Extract the URL from the command
    url := strings.TrimSpace(strings.TrimPrefix(content, "!play "))
 
    // Download the audio file from the URL
    audioFile, err := downloadAudioFile(url)
    if err != nil {
        return err
    }
 
    // Encode the audio file into a DCA format
    dcaFile, err := encodeAudioFile(audioFile)
    if err != nil {
        return err
    }
 
    // Play the audio file in the voice channel
    done := make(chan error)
    stream := dca.NewStream(bytes.NewReader(dcaFile), bot.Voice, done)
    defer bot.Voice.Speaking(false)
    bot.Voice.Speaking(true)
    defer close(done)
    for {
        select {
        case err := <-done:
            if err != nil && err != io.EOF {
                return err
            }
            return nil
        default:
            _, err := io.Copy(ffmpeg.NewEncoder(bot.Voice), stream)
            if err != nil && err != io.EOF {
                return err
            }
        }
    }
}
 
func downloadAudioFile(url string) ([]byte, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
 
    audioFile, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
 
    return audioFile, nil
}
 
func encodeAudioFile(audioFile []byte) ([]byte, error) {
    cmd := exec.Command("ffmpeg", "-i", "-", "-f", "dca", "-")
    stdin, err := cmd.StdinPipe()
    if err != nil {
        return nil, err
    }
 
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return nil, err
    }
 
    err = cmd.Start()
    if err != nil {
        return nil, err
    }
 
    _, err = stdin.Write(audioFile)
    if err != nil {
        return nil, err
    }
 
    stdin.Close()
 
    dcaFile, err := ioutil.ReadAll(stdout)
    if err != nil {
        return nil, err
    }
 
    err = cmd.Wait()
    if err != nil {
        return nil, err
    }
 
    return dcaFile, nil
}

func main() {
    bot := &Bot{}
 
    discordToken := "YOUR_DISCORD_TOKEN"
    discord, err := discordgo.New("Bot " + discordToken)
    if err != nil {
        fmt.Println("Error creating Discord session:", err)
        return
    }
 
    bot.Session = discord
    bot.Session.AddHandler(bot.onMessageCreate)
 
    err = bot.Session.Open()
    if err != nil {
        fmt.Println("Error opening Discord session:", err)
        return
    }
 
    fmt.Println("Bot is now running. Press CTRL-C to exit.")
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sigChan
 
    bot.Session.Close()
}