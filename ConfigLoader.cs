using System.Text.Json;

namespace CornealiusEyeworth;

class ConfigLoader
{
    public Config Load()
    {
        var configPath = Path.Combine(AppContext.BaseDirectory, "config.json");
        return JsonSerializer.Deserialize<Config>(File.ReadAllText(configPath))!;
    }
}
