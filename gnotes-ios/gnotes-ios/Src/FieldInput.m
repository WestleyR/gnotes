//
//  FieldInput.m
//  gnotes-ios
//
//  Created by Westley Rose on 2/25/23.
//

#import "FieldInput.h"

@implementation FieldInput

// TODO: unused


-(BOOL)textFieldShouldBeginEditing:(UITextField *)textField {
    NSLog(@"%s", __func__);

    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(keyboardDidShow:) name:UIKeyboardDidShowNotification object:nil];
    return YES;
}


- (BOOL)textFieldShouldEndEditing:(UITextField *)textField {
    NSLog(@"%s", __func__);

    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(keyboardDidHide:) name:UIKeyboardDidHideNotification object:nil];

    [self endEditing:YES];
    return YES;
}


- (void)keyboardDidShow:(NSNotification *)notification {
    NSLog(@"%s", __func__);

    // Assign new frame to your view
    [self setFrame:CGRectMake(0,-110,320,460)]; //here taken -110 for example i.e. your view will be scrolled to -110. change its value according to your requirement.

}

-(void)keyboardDidHide:(NSNotification *)notification {
    NSLog(@"%s", __func__);

    [self setFrame:CGRectMake(0,0,320,460)];
}


/*
// Only override drawRect: if you perform custom drawing.
// An empty implementation adversely affects performance during animation.
- (void)drawRect:(CGRect)rect {
    // Drawing code
}
*/

@end
