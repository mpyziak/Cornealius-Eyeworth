namespace CornealiusEyeworth;

interface IConfigRepository
{
    Config Load();
    void Save(Config config);
}
