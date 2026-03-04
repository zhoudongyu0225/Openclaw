package main

// 战斗逻辑
type Battle struct {
    Wave int
    Enemies []Enemy
}

type Enemy struct {
    Type string
    HP int
    Speed int
}

func (b *Battle) StartWave() {
    b.Wave++
    count := 5 + b.Wave * 2
    for i := 0; i < count; i++ {
        b.Enemies = append(b.Enemies, Enemy{
            Type: "droid",
            HP: 100 + b.Wave * 10,
            Speed: 2,
        })
    }
}
