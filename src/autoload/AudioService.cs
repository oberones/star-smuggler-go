using Godot;

namespace StarSmugglerGo.Autoload;

public partial class AudioService : Node
{
    private readonly System.Collections.Generic.Dictionary<string, string> _musicTracks = new()
    {
        ["singularity"] = "res://assets/audio/music/singularity.mp3",
        ["world_default"] = "res://assets/audio/music/world_default.mp3",
    };

    private readonly System.Collections.Generic.Dictionary<string, string> _sfxTracks = new()
    {
        ["click"] = "res://assets/audio/sfx/click.wav",
    };

    private AudioStreamPlayer? _musicPlayer;
    private AudioStreamPlayer? _sfxPlayer;

    public string? CurrentMusicTrackId { get; private set; }

    public override void _Ready()
    {
        _musicPlayer = GetNodeOrNull<AudioStreamPlayer>("MusicPlayer");
        _sfxPlayer = GetNodeOrNull<AudioStreamPlayer>("SfxPlayer");
    }

    public void PlayMusic(string trackId)
    {
        if (CurrentMusicTrackId == trackId)
        {
            return;
        }

        CurrentMusicTrackId = trackId;

        if (_musicPlayer is null)
        {
            return;
        }

        if (!_musicTracks.TryGetValue(trackId, out string? path))
        {
            GD.PushWarning($"AudioService could not resolve music track '{trackId}'.");
            return;
        }

        AudioStream? stream = ResourceLoader.Load<AudioStream>(path);
        if (stream is null)
        {
            GD.PushWarning($"AudioService failed to load music stream '{path}'.");
            return;
        }

        _musicPlayer.Stream = stream;
        _musicPlayer.Play();
    }

    public void PlaySfx(string sfxId)
    {
        if (_sfxPlayer is null)
        {
            return;
        }

        if (!_sfxTracks.TryGetValue(sfxId, out string? path))
        {
            GD.PushWarning($"AudioService could not resolve SFX '{sfxId}'.");
            return;
        }

        AudioStream? stream = ResourceLoader.Load<AudioStream>(path);
        if (stream is null)
        {
            GD.PushWarning($"AudioService failed to load SFX stream '{path}'.");
            return;
        }

        _sfxPlayer.Stream = stream;
        _sfxPlayer.Play();
    }
}
