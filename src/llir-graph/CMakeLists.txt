file(GLOB_RECURSE LLIR_GRAPH_SOURCE_FILES 
		${CMAKE_CURRENT_SOURCE_DIR}/*.cpp)

file(GLOB_RECURSE LLIR_GRAPH_HEADER_FILES 
		${CMAKE_CURRENT_SOURCE_DIR}/*.hpp)

add_executable(${PROJECT_NAME} ${LLIR_GRAPH_SOURCE_FILES} ${LLIR_GRAPH_HEADER_FILES})
set_compile_options(${PROJECT_NAME})
target_include_directories(${PROJECT_NAME} PRIVATE ./ ${LLVM_INCLUDE_DIRS})
# separate_arguments(LLVM_DEFINITIONS_LIST NATIVE_COMMAND ${LLVM_DEFINITIONS})
# target_compile_definitions(${PROJECT_NAME} PRIVATE ${LLVM_DEFINITIONS_LIST})
llvm_map_components_to_libnames(llvm_libs support core irreader analysis)
target_link_libraries(${PROJECT_NAME} PRIVATE fmt::fmt-header-only ${llvm_libs} stdc++)
