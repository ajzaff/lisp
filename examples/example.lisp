(doc fact returns the factorial of i

    (Example
        (eq (fact 8) 40320))

    (SpecialCases
        (eq (fact 0) 1)
        (eq (fact 1) 1))
)
(def fact (i) 
    (if (eq i 0) 
        1
    (elif (eq i 1)
        1
    (else
        (mul i (fact (sub i 1))
    )))))