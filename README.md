# Gosh 🐚

Um mini shell escrito em Go — simples, poderoso, extensível.

> Futuro lar do `oh-my-gosh` ✨

---

## 🚀 O que é?

**Gosh** é um shell minimalista implementado em Go, com objetivo educacional e prático, projetado para te dar controle total sobre o terminal enquanto explora os bastidores de sistemas operacionais, processos e sinais.

---

## ✅ Funcionalidades Atuais

- Leitura de comandos linha a linha (REPL)
- Execução de comandos externos simples (`ls`, `echo`, etc)
- Arquitetura modular (parser, executor, etc)

---

## 📦 Estrutura

```bash
gosh/
├── cmd/
│   └── gosh.go           # entrada da aplicação
├── internal/
│   ├── parser/           # parsing simples (tokens)
│   └── executor/         # execução de comandos externos
├── go.mod
└── README.md
```

---

## ▶️ Como rodar

```bash
git clone https://github.com/seu-usuario/gosh.git
cd gosh
go run ./cmd/gosh.go
```

---

## 💡 Roadmap

- [x] Execução de comandos externos
- [ ] Comandos internos (`cd`, `exit`)
- [ ] Pipes (`|`)
- [ ] Redirecionamento (`>`, `<`, `>>`)
- [ ] Execução em background (`&`)
- [ ] Sinais (`SIGINT`, `SIGTSTP`, etc)
- [ ] Job Control (`jobs`, `fg`, `bg`, `kill`)
- [ ] Histórico de comandos (`~/.gosh_history`)
- [ ] Prompt customizável (usuário, host, path)
- [ ] Autocompletar básico (TAB)
- [ ] Modo script `.goshrc`

---

## 🌈 Futuro: oh-my-gosh

Uma suíte de extensões, temas e plugins inspirada no `oh-my-zsh`.

---

## 📚 Feito para aprender

Este projeto cobre:

- Processos e fork/exec em Unix
- Sinais e controle de jobs
- Parsers e tokenização
- Goroutines, canais, sistema de arquivos
- Go idiomático e modular

---

## 📜 Licença

MIT © Felipe Fernandes
