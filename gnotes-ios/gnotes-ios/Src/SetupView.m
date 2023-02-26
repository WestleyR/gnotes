//
//  SetupView.m
//  gnotes-ios
//
//  Created by Westley Rose on 2/25/23.
//

// TODO: REMOVE AS THIS IS REPLACES BY SetupViewController.swift

#import "SetupView.h"

@interface SetupView ()

@end

@implementation SetupView

/*
- (void)viewDidLoad {
    [super viewDidLoad];
    // Do any additional setup after loading the view.

    self.buttonDone.enabled = NO;
}

- (void)canEnableDone {
    int entered = 0;

    if (![self.fieldAccessKey.text isEqual: @""]) {
        entered++;
    }
    if (![self.fieldSecretKey.text isEqual: @""]) {
        entered++;
    }
    if (![self.fieldAccountID.text isEqual: @""]) {
        entered++;
    }
    if (![self.fieldCryptKey.text isEqual: @""]) {
        entered++;
    }

    if (entered == 4) {
        self.buttonDone.enabled = YES;
    } else {
        [self.buttonDone setEnabled:NO];
    }
}

- (IBAction)buttonDone:(id)sender {
    NSLog(@"%s", __func__);

    [self dismissViewControllerAnimated:YES completion:nil];
    [self canEnableDone];
}

- (IBAction)fieldAccessKey:(id)sender {
    NSLog(@"%s", __func__);

    [self.view endEditing:YES];
    [self canEnableDone];
}

- (IBAction)fieldSecretKey:(id)sender {
    NSLog(@"%s", __func__);

    [self.view endEditing:YES];
    [self canEnableDone];
}

- (IBAction)fieldAccountID:(id)sender {
    NSLog(@"%s", __func__);

    [self.view endEditing:YES];
    [self canEnableDone];
}

- (IBAction)fieldCryptKey:(id)sender {
    NSLog(@"%s", __func__);

    [self.view endEditing:YES];
    [self canEnableDone];
}


#pragma mark - Navigation

// In a storyboard-based application, you will often want to do a little preparation before navigation
- (void)prepareForSegue:(UIStoryboardSegue *)segue sender:(id)sender {
    // Get the new view controller using [segue destinationViewController].
    // Pass the selected object to the new view controller.

    [self canEnableDone];
}

*/
@end
