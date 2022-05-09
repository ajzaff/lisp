('doc fact returns the factorial of ``i''

Example:
    (= (fact 8) 40320)

Special cases:
    (= (fact 0) 1)
    (= (fact 1) 1)
)
(def fact (i) 
    (if (= i 0) 
        1
    (elif (= i 1)
        1
    (else
        (* i (fact (- i 1))
    )))))