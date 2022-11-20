#pragma once

#include <filesystem>

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
            if (validate(args.inputFilePath().empty(), "input path not given"))
                validate(!std::filesystem::exists(args.inputFilePath()), "no file exists at given path");
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
