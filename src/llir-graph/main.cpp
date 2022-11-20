#include "cli/cli.hpp"
#include "cli/cli_check.hpp"

#include <fmt/core.h>
#include <iostream>

int main(int argc, char* argv[])
{
    CLIArgs args(argc, argv);

    {
        auto errors = llir_graph::cli::CLICheck(args).dumpErrors();
        if (!errors.empty()) {
            for (const auto& error : errors) {
                std::clog << fmt::format("args parsing error: {}\n", error);
            }
        }
    }

    return 0;
}
