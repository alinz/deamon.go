# DEAMON.go

I found myself keep writing the same signal intrupt code over and over again, so I created this. 

## INSTALLATION

do the either of options:

- you can even copy the code in `deamon.go` to your project directly
- or run `go get github.com/alinz/deamon.go` inside your go project


## USAGE

simply in your main program use it as follows;


```golang

func main() {
    ctx := context.Background()

    err := deamon.Summoning(ctx, func (ctx context.Context, summonType deamon.SummonType) error {
        switch summonType {
            case deamon.Call:
                // only happens once at the start, go point for initializing
                // all long lived objects here
            case deamon.Recall:
                // program requested to be reloaded
                // might be a good point to reload config files here
            case deamon.Kill:
                // program requested to be killed
                // do the neccessary clean up here
        }

        return nil
    })

    if err != nil {
        log.Println(err)
        os.Exit(1)
    }
}
```