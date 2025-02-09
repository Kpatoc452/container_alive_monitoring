package main

import "github.com/Kpatoc452/VK_exams/WorkerPool-Task/logger"


type OptionWP struct{
    Max int
    Logger logger.Logger
}


type Option func(o *OptionWP)


func NewOptionWP(opts ...Option) OptionWP{
    opt := OptionWP{
        Max: 10,
        Logger: logger.New(),
    }

    for _, o := range opts{
        o(&opt)
    }

    return opt
}