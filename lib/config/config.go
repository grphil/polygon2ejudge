package config

// Global configs

const JUDGES_DIR = "/Users/grphil/Documents/ejudge/judges"
const GVALUER_LOCATION = "/Users/grphil/Documents/ejudge/judges/001501/problems/gvaluer"
const EJUDGE_URL = "https://ejudge.algocode.ru/"

// Problem configs

const CREATE_STATEMENTS = true
const COMPILE_ALL_SOLUTIONS = false // If false, all solutions are imported in solutions1 folder. They still can be used for autosubmit
const COMPILE_MAIN_SOLUTION = false // If false, no solution is set for the problem
const TEXTAREA_INPUT = true
const ENABLE_CUSTOM_RUN = true
const GENERIC_PARENT = true           // If true, if problem has abstract problem with name Generic, it is used as parent
const NOLINT_STRING = "secret"        // This string will be added as comment to all solutions with nolint submit
const FULL_REPORT_ONLY_SAMPLES = true // If true, full output will be shown only on samples

var LANG_IDS = map[string][]int{
	"py":   {23, 64}, // python3 and pypy3
	"cpp":  {3, 52},  // g++ and clang
	"java": {18},     // java
	"pas":  {1},      // Free Pascal
}
