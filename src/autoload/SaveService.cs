using Godot;

namespace StarSmugglerGo.Autoload;

public partial class SaveService : Node
{
    public const int CurrentSaveVersion = 1;

    public string SavePath => "user://save.json";
}
