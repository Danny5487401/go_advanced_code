#include "number.h"

int number_add_mod(int a, int b, int mod) {
    return (a+b)%mod;
}


// gcc -shared -o libnumber.so number.c

// 最好放在动态库位置
// 否则报错
//   Referenced from: /private/var/folders/sk/m49vysmj3ss_y50cv9cvcn800000gn/T/___go_build_main_go__2_
//    Reason: tried: 'libnumber.so' (no such file), '/usr/local/lib/libnumber.so' (no such file), '/usr/lib/libnumber.so' (no such file), '/Users/python/Desktop/github.com/Danny5487401/go_advanced_code/libnumber.so' (no such file), '/usr/local/lib/libnumber.so' (no such file), '/usr/lib/libnumber.so' (no such file)