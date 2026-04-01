using Godot;

namespace StarSmugglerGo.Presentation;

public partial class AppRoot : Control
{
    public override void _Ready()
    {
        EnsureDefaultAction("menu_confirm", Key.Enter, Key.Space);
        EnsureDefaultAction("menu_back", Key.Escape, Key.Backspace);
        EnsureDefaultAction("menu_up", Key.Up, Key.W);
        EnsureDefaultAction("menu_down", Key.Down, Key.S);
        EnsureDefaultAction("menu_left", Key.Left, Key.A);
        EnsureDefaultAction("menu_right", Key.Right, Key.D);
    }

    private static void EnsureDefaultAction(string actionName, params Key[] keys)
    {
        if (!InputMap.HasAction(actionName))
        {
            InputMap.AddAction(actionName);
        }

        foreach (Key key in keys)
        {
            var inputEvent = new InputEventKey
            {
                Keycode = key,
                PhysicalKeycode = key,
            };

            if (!InputMap.ActionHasEvent(actionName, inputEvent))
            {
                InputMap.ActionAddEvent(actionName, inputEvent);
            }
        }
    }
}
