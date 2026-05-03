namespace StarSmugglerGo.Services;

public sealed class UpgradePurchaseResult
{
    public bool Succeeded { get; init; }
    public string Message { get; init; } = string.Empty;

    public static UpgradePurchaseResult Success(string message) => new()
    {
        Succeeded = true,
        Message = message,
    };

    public static UpgradePurchaseResult Failure(string message) => new()
    {
        Succeeded = false,
        Message = message,
    };
}
