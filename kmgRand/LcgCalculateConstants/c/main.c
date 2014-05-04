#include "rand-lcg.h"
#include "rand-primegen.h" /* DJB's prime factoring code */
#include "stdio.h"
/*
 come from https://github.com/robertdavidgraham/masscan
*/
int main(int argc, char *argv[]){
  uint64_t m=0;
  uint64_t a=0;
  uint64_t c=0;
  if (argc==2){
    sscanf(argv[1],"%llu",&m);
  }else if(argc==3){
    sscanf(argv[2],"%llu",&c);
  }else{
    printf("usage: %s [m(the range of lcg)] [c]\n",argv[0]);
    return -1;
  }
  sscanf(argv[1],"%llu",&m);
  puts("random_value = (index * a + c) % range;\n");
  lcg_calculate_constants(m,&a,&c,1);
  return 0;
}

