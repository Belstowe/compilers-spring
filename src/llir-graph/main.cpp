#include "cli/cli.hpp"
#include "cli/cli_check.hpp"
#include "graph/call.hpp"
#include "graph/def_use.hpp"

#include <fmt/core.h>
#include <iostream>

#include <llvm/Analysis/CallGraph.h>
#include <llvm/Support/SourceMgr.h>
#include <llvm/IR/Module.h>
#include <llvm/IR/LLVMContext.h>
#include <llvm/IRReader/IRReader.h>

int main(int argc, char *argv[])
{
    CLIArgs args(argc, argv);

    {
        auto errors = llir_graph::cli::CLICheck(args).dumpErrors();
        if (!errors.empty())
        {
            for (const auto &error : errors)
            {
                std::clog << fmt::format("args parsing error: {}\n", error);
            }
            return 1;
        }
    }

    llvm::SMDiagnostic err;
    llvm::LLVMContext context;
    auto irModule = llvm::parseIRFile(args.inputFilePath(), err, context);
    //llvm::CallGraph graph(*(irModule.get()));

    if (args.defUseGraph()) {
        std::cout << fmt::format("digraph {{\n{} }}\n", llir_graph::graph::defUseGraphToString(irModule.get()));
    } else if (args.callGraph()) {
        std::cout << fmt::format("digraph {{\n{} }}\n", llir_graph::graph::callGraphToString(irModule.get()));
    }
    return 0;
}
