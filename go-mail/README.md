# go-mail

## 怎么发

很简单：
```golang
to0 := []string{
    "xx@xx.com",
    "yy@yy.com",
}

to1 := []string{
    "zz@zz.com",
}

NewSender().Login("sender@sender.com","sender_pwd","smtp.sender.com","smtp.sender.com:587").Send(subject0,message0,to0...).Send(subject1,message1,to1...).Done()
```