#pragma once

#include <llvm/Support/CommandLine.h>
#include <string>

namespace llir_graph::cli
{
    class CLIArgs
    {
    private:
        llvm::cl::opt<std::string> outputFilePathArg_;
        llvm::cl::opt<std::string> inputFilePathArg_;
        llvm::cl::opt<bool> callGraphArg_;
        llvm::cl::opt<bool> defUseGraphArg_;

    public:
        CLIArgs(int argc, char *argv[]) : outputFilePathArg_("output"),
                                          inputFilePathArg_("input"),
                                          callGraphArg_("dot-callgraph"),
                                          defUseGraphArg_("dot-def-use")
        {
            llvm::cl::ParseCommandLineOptions(argc, argv);
        }

        const std::string &outputFilePath() const
        {
            return outputFilePathArg_.getValue();
        }

        const std::string &inputFilePath() const
        {
            return inputFilePathArg_.getValue();
        }

        bool callGraph() const
        {
            return callGraphArg_.getValue();
        }

        bool defUseGraph() const
        {
            return defUseGraphArg_.getValue();
        }
    };
}
