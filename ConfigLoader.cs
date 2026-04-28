using System.Text.Json;

namespace CornealiusEyeworth;

class ConfigLoader
{
    public Config Load()
    {
        var exeDir = Path.GetDirectoryName(Environment.ProcessPath)!;
        var configPath = Path.Combine(exeDir, "config.json");
        return JsonSerializer.Deserialize<Config>(File.ReadAllText(configPath))!;
    }
}