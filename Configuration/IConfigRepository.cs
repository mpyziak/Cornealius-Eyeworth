namespace CornealiusEyeworth.Configuration;

/// <summary>
/// Abstracts persistence of <see cref="Config"/> so that the storage
/// mechanism (JSON, registry, cloud, …) can be swapped without touching
/// any consumer.
/// </summary>
internal interface IConfigRepository
{
    /// <summary>Reads and returns the current configuration.</summary>
    Config Load();

    /// <summary>Persists <paramref name="config"/> to the backing store.</summary>
    void Save(Config config);
}
