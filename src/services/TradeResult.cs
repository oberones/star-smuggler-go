namespace StarSmugglerGo.Services;

public sealed class TradeResult
{
    public bool Succeeded { get; init; }
    public string Message { get; init; } = string.Empty;

    public static TradeResult Success(string message)
    {
        return new TradeResult
        {
            Succeeded = true,
            Message = message,
        };
    }

    public static TradeResult Failure(string message)
    {
        return new TradeResult
        {
            Succeeded = false,
            Message = message,
        };
    }
}
