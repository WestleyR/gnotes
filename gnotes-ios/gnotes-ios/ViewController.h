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


@property (weak, nonatomic) IBOutlet UITableView *tableViewNotes;

@end

