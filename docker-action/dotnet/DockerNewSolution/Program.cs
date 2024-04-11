using CliWrap;
using CommandLine;

Parser.Default.ParseArguments<Options>(args)
    .WithParsed<Options>(o =>
    {
        var options = ""; 
        var cmd = Cli.Wrap("dotnet")
            .WithArguments(args => args
                .Add("new")
                .Add("sln")
                .Add(options)
            );
    });

internal class Options
{
    [Option('s', "solutionName", Required = true, HelpText = "Specify the solution name.")]
    public string SolutionName { get; set; } = string.Empty;

    [Option('o', "outputFolder", Required = false, HelpText = "Specify the output folder.")]
    public string OutputFolder { get; set; } = string.Empty;
}