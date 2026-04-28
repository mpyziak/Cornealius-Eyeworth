using System.Text.Json;

namespace CornealiusEyeworth;

class JsonConfigRepository : IConfigRepository
{
    private static readonly JsonSerializerOptions WriteOptions = new() { WriteIndented = true };

    private readonly string _configPath;

    public JsonConfigRepository()
    {
        var exeDir = Path.GetDirectoryName(Environment.ProcessPath)!;
        _configPath = Path.Combine(exeDir, "config.json");
    }

    public Config Load()
    {
        var json = File.ReadAllText(_configPath);
        return JsonSerializer.Deserialize<Config>(json)!;
    }

    public void Save(Config config)
    {
        var json = JsonSerializer.Serialize(config, WriteOptions);
        File.WriteAllText(_configPath, json);
    }
}
