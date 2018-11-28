# monitor

Enable the monitor by using `-m` on the command line.

## Arguments

The arguments for *address* and *value* can be specified as a hexadecimal value without a prefix or with a `$` prefix. Decimal values can be specified by using a `+` prefix.

Examples:
```
p 1234 ff
p 1234 $ff
p 1234 +255
```

## Commands

### b *address* {on|off}

Sets a **breakpoint** at *address* when using `on` and clears a breakpoint at *address* when using `off`. The CPU stops before executing *address*.

### d [*start-address* [*end-address*]]

**Disassemble** code from *start-address* to *end-address* inclusive. If *end-address* is not specified, disassemble an amount that can fit on a screen. If *start-address* is not specified, use the current program counter as the *start-address*.

### f *start-address* *end-address* *value*

**Fill** memory with *value* from *start-address* to *end-address* inclusive.

### g [*address*]

**Go** to *address* and start execution of the CPU there. If *address* is not specified, use the current value of the program counter.

### h

**Halt** execution of the CPU.

### m [*start-address* [*end-address*]]

Dump **memory** contents to the screen from *start-address* to *end-address* inclusive. If *end-address* is not specified, show a full memory page. If *start-address* is not specified, continue the dump from the last command.

### n

Disassemble the next instruction to execute.

### p *address*

**Peek** at the memory contents at *address*. The value is displayed in the form of `$00 +000` with the hexadecimal value listed first followed by the decimal value.

### p *address* *value*

**Poke** the memory at *address* with *value*.

### r

Display the contents of the CPU **registers**.

### r *name*

Display the value for the **register** with *name*.

### r *name* *value*

Set the *value* for **register** with *name*.

### s

**Step** through by executing the next instruction and then halting the CPU.

### so

Save the current machine **state out** to disk .

### si

Load the current machine **state in** from disk.

### t

Toggle **tracing** of instructions executed by the CPU.

### q[uit]

Quit to the operating system.
