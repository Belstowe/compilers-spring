#include <stdio.h>

int factorial(int n)
{
    int r = 1;
    for (; n > 1; r *= n--);
    return r;
}

int main()
{
    int n;
    int f;

    printf("Enter an integer to find its factorial\n");
    scanf("%d", &n);

    if (n < 0)
    {
        printf("Factorial of negative integers isn't defined.\n");
        return 1;
    }

    f = factorial(n);
    printf("%d! = %d\n", n, f);

    return 0;
}
