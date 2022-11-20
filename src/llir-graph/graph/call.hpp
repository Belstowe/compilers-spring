#pragma once

#include <string>
#include <unordered_set>

#include <fmt/core.h>
#include <llvm/Analysis/CallGraph.h>
#include <llvm/IR/Instructions.h>
#include <llvm/IR/InstrTypes.h>
#include <llvm/IR/InstIterator.h>
#include <llvm/IR/Module.h>

namespace llir_graph::graph
{
    std::string callGraphToString(llvm::Module *module)
    {
        std::string str;
        std::unordered_set<std::string> alreadyCalledFunctions;

        for (const auto& func : module->functions())
        {
            for (auto instr = llvm::inst_begin(func); instr != llvm::inst_end(func); instr++)
            {
                if (llvm::isa<llvm::CallInst>(*instr))
                {
                    llvm::StringRef name = llvm::cast<llvm::CallInst>(*instr).getCalledFunction()->getName();
                    if (alreadyCalledFunctions.find(std::string(name.data())) == alreadyCalledFunctions.end())
                    {
                        str += fmt::format("{} -> {};\n", func.getName().data(), name.data());
                        alreadyCalledFunctions.insert(std::string(name.data()));
                    }
                }
            }
            alreadyCalledFunctions.clear();
        }

        return str;
    }
}
