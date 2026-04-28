using System.Text.Json;

namespace CornealiusEyeworth.Configuration;

/// <summary>
/// Persists <see cref="Config"/> as a human-readable JSON file located
/// next to the application executable.
/// </summary>
internal class JsonConfigRepository : IConfigRepository
{
    private static readonly JsonSerializerOptions WriteOptions = new() { WriteIndented = true };

    private readonly string _configPath;

    /// <summary>
    /// Resolves the config file path relative to the running executable's directory.
    /// </summary>
    public JsonConfigRepository()
    {
        var exeDir = Path.GetDirectoryName(Environment.ProcessPath)!;
        _configPath = Path.Combine(exeDir, "config.json");
    }

    /// <inheritdoc/>
    public Config Load()
    {
        var json = File.ReadAllText(_configPath);
        return JsonSerializer.Deserialize<Config>(json)!;
    }

    /// <inheritdoc/>
    public void Save(Config config)
    {
        var json = JsonSerializer.Serialize(config, WriteOptions);
        File.WriteAllText(_configPath, json);
    }
}
