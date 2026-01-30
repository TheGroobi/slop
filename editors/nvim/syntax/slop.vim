if exists("b:current_syntax")
    finish
endif

" Comments
syn match slopComment "#.*$"

" Directives: config, var, run (before ::)
syn match slopDirective "\<\(config\|var\|run\)\>\ze\s*::"

" Double colon separator
syn match slopSeparator "::"

" Nested key depths: key1.key2.key3
" Depth 1: first identifier after ::  (e.g. db)
syn match slopKey1 "\(::\s*\)\@<=\w\+" nextgroup=slopDot1

" Actions (after run::): seed, migrate, backup, dump — defined after slopKey1 for priority
syn match slopAction "\(run\s*::\)\@<=\s*\<\(seed\|migrate\|backup\|dump\)\>"
" Depth 2: after first dot            (e.g. type)
syn match slopDot1 "\." contained nextgroup=slopKey2
syn match slopKey2 "\w\+" contained nextgroup=slopDot2
" Depth 3: after second dot           (e.g. foo)
syn match slopDot2 "\." contained nextgroup=slopKey3
syn match slopKey3 "\w\+" contained

" Brackets
syn match slopBracket "\[\|\]"

" Strings
syn region slopString start='"' end='"'

" Variable references
syn match slopVarRef "\$\w\+\(\.\w\+\)*"

" Highlight links
hi link slopComment    Comment
hi link slopDirective  Keyword
hi link slopSeparator  Operator
hi link slopAction     Function
hi link slopKey1       Identifier
hi link slopDot1       Operator
hi link slopKey2       Type
hi link slopDot2       Operator
hi link slopKey3       Constant
hi link slopBracket    Delimiter
hi link slopString     String
hi link slopVarRef     Special

let b:current_syntax = "slop"
