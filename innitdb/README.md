# innitdb

(?x (1 2 3)) => 

(?x ?y) where ?y = (1 2 3):


SELECT ?x
    FROM t
    WHERE
    (?x (1 2 3)) in t