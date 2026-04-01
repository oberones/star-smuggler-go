using Godot;

namespace StarSmugglerGo.Autoload;

public partial class AudioService : Node
{
    public string? CurrentMusicTrackId { get; private set; }

    public void PlayMusic(string trackId)
    {
        CurrentMusicTrackId = trackId;
    }
}
