00000000  FC                cld
00000001  4883E4F0          and rsp,0xfffffffffffffff0
00000005  E8C0000000        call 0xca
0000000A  4151              push r9
0000000C  4150              push r8
0000000E  52                push rdx
0000000F  51                push rcx
00000010  56                push rsi
00000011  4831D2            xor rdx,rdx
00000014  65488B5260        mov rdx,[gs:rdx+0x60]
00000019  488B5218          mov rdx,[rdx+0x18]
0000001D  488B5220          mov rdx,[rdx+0x20]
00000021  488B7250          mov rsi,[rdx+0x50]
00000025  480FB74A4A        movzx rcx,word [rdx+0x4a]
0000002A  4D31C9            xor r9,r9
0000002D  4831C0            xor rax,rax
00000030  AC                lodsb
00000031  3C61              cmp al,0x61
00000033  7C02              jl 0x37
00000035  2C20              sub al,0x20
00000037  41C1C90D          ror r9d,byte 0xd
0000003B  4101C1            add r9d,eax
0000003E  E2ED              loop 0x2d
00000040  52                push rdx
00000041  4151              push r9
00000043  488B5220          mov rdx,[rdx+0x20]
00000047  8B423C            mov eax,[rdx+0x3c]
0000004A  4801D0            add rax,rdx
0000004D  8B8088000000      mov eax,[rax+0x88]
00000053  4885C0            test rax,rax
00000056  7467              jz 0xbf
00000058  4801D0            add rax,rdx
0000005B  50                push rax
0000005C  8B4818            mov ecx,[rax+0x18]
0000005F  448B4020          mov r8d,[rax+0x20]
00000063  4901D0            add r8,rdx
00000066  E356              jrcxz 0xbe
00000068  48FFC9            dec rcx
0000006B  418B3488          mov esi,[r8+rcx*4]
0000006F  4801D6            add rsi,rdx
00000072  4D31C9            xor r9,r9
00000075  4831C0            xor rax,rax
00000078  AC                lodsb
00000079  41C1C90D          ror r9d,byte 0xd
0000007D  4101C1            add r9d,eax
00000080  38E0              cmp al,ah
00000082  75F1              jnz 0x75
00000084  4C034C2408        add r9,[rsp+0x8]
00000089  4539D1            cmp r9d,r10d
0000008C  75D8              jnz 0x66
0000008E  58                pop rax
0000008F  448B4024          mov r8d,[rax+0x24]
00000093  4901D0            add r8,rdx
00000096  66418B0C48        mov cx,[r8+rcx*2]
0000009B  448B401C          mov r8d,[rax+0x1c]
0000009F  4901D0            add r8,rdx
000000A2  418B0488          mov eax,[r8+rcx*4]
000000A6  4801D0            add rax,rdx
000000A9  4158              pop r8
000000AB  4158              pop r8
000000AD  5E                pop rsi
000000AE  59                pop rcx
000000AF  5A                pop rdx
000000B0  4158              pop r8
000000B2  4159              pop r9
000000B4  415A              pop r10
000000B6  4883EC20          sub rsp,0x20
000000BA  4152              push r10
000000BC  FF                db 0xff
000000BD  E058              loopne 0x117
000000BF  4159              pop r9
000000C1  5A                pop rdx
000000C2  488B12            mov rdx,[rdx]
000000C5  E957FFFFFF        jmp 0x21
000000CA  5D                pop rbp
000000CB  48BA010000000000  mov rdx,0x1
         -0000
000000D5  488D8D01010000    lea rcx,[rbp+0x101]
000000DC  41BA318B6F87      mov r10d,0x876f8b31
000000E2  FF                db 0xff
000000E3  D5BBE01D2A0A41BA  o64 pavgb mm3,[rel 0xffffffffba410b15]
000000EB  A6                cmpsb
000000EC  95                xchg eax,ebp
000000ED  BD9DFFD548        mov ebp,0x48d5ff9d
000000F2  83C428            add esp,0x28
000000F5  3C06              cmp al,0x6
000000F7  7C0A              jl 0x103
000000F9  80FBE0            cmp bl,0xe0
000000FC  7505              jnz 0x103
000000FE  BB4713726F        mov ebx,0x6f721347
00000103  6A00              push qword 0x0
00000105  59                pop rcx
00000106  4189DA            mov r10d,ebx
00000109  FF                db 0xff
0000010A  D5                db 0xd5
0000010B  63616C            movsxd esp,dword [rcx+0x6c]
0000010E  632E              movsxd ebp,dword [rsi]
00000110  657865            gs js 0x178
00000113  00                db 0x00
