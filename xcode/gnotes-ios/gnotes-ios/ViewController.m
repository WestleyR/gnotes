//
//  ViewController.m
//  gnotes-ios
//
//  Created by Westley Rose on 3/6/22.
//

#import "ViewController.h"

#include "gnotes.h"

@interface ViewController ()

@end

@implementation ViewController

- (IBAction)button:(id)sender {
    NSLog(@"HELLO");

    char* foo = List("");
    NSLog(@"LISTT RESP: %s", foo);
    free(foo);

}

- (void)viewDidLoad {
    [super viewDidLoad];
    // Do any additional setup after loading the view.
}


@end
