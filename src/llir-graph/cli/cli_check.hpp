#pragma once

#include "cli.hpp"

using llir_graph::cli::CLIArgs;

namespace llir_graph::cli
{
    class CLICheck
    {
    private:
        std::vector<std::string> cliErrors_;

        bool validate(bool cond, const char *thenMessage)
        {
            if (cond)
            {
                cliErrors_.push_back(std::string(thenMessage));
            }
            return !cond;
        }

    public:
        CLICheck(const CLIArgs &args)
            : cliErrors_()
        {
            validate(args.inputFilePath().empty(), "input path not given");
            // validate(args.outputFilePath().empty(), "output path not given");
            int graphFlagCount = args.callGraph() + args.defUseGraph();
            validate(graphFlagCount > 1, "more than one graph flag specified");
            validate(graphFlagCount < 1, "no graph flags specified");
        }

        const std::vector<std::string> &dumpErrors() const
        {
            return cliErrors_;
        }
    };
}
