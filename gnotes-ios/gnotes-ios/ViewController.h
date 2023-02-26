//
//  ViewController.h
//  gnotes-ios
//
//  Created by Westley Rose on 2/23/23.
//

#import <UIKit/UIKit.h>

#import "SetupView.h"

@interface ViewController : UIViewController {
    NSMutableArray *noteTitles;
}

// REMOVE:
@property (weak, nonatomic) IBOutlet UILabel *labelResponse;
@property (weak, nonatomic) IBOutlet UIActivityIndicatorView *loadingWheelMiddle;



@property (weak, nonatomic) IBOutlet UITableView *tableViewNotes;

@end

