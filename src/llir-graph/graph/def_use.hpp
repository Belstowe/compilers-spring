#pragma once

#include <string>

#include <fmt/core.h>
#include <llvm/Analysis/CallGraph.h>
#include <llvm/IR/Instructions.h>
#include <llvm/IR/InstrTypes.h>
#include <llvm/IR/User.h>

namespace llir_graph::graph
{
    std::string defUseGraphToString(llvm::Module *module)
    {
        std::string str;
        llvm::raw_string_ostream stream(str);

        for (const auto& func : module->functions())
        {
            for (auto instr = llvm::inst_begin(func); instr != llvm::inst_end(func); instr++) {
                for (const auto &operand : instr->operands())
                {
                    llvm::Value *v = operand.get();
                    if (llvm::dyn_cast<llvm::Instruction>(v))
                    {
                        stream << "\"" << *instr << "\" -> \"" << *llvm::dyn_cast<llvm::Instruction>(v) << "\";\n";
                    }
                }
            }
        }

        return str;
    }
}
